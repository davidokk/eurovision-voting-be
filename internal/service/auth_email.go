package service

import (
	"context"
	"crypto/rand"
	"errors"
	"eurovision-voting/internal/domain"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

const (
	authPurposeSignup        = "signup"
	authPurposeEmailBind     = "email_bind"
	authPurposePasswordReset = "password_reset"

	authCodeTTL        = 15 * time.Minute
	emailRateLimitMax  = 5
	emailRateLimitWindow = time.Hour
)

func normalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return "", ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrInvalidEmail
	}
	return email, nil
}

func (s *Service) IsTelegramLinked(u *domain.User) bool {
	return u.TelegramID != nil && *u.TelegramID > 0
}

func (s *Service) IsEmailVerified(u *domain.User) bool {
	if s.IsTelegramLinked(u) {
		return true
	}
	return u.EmailVerifiedAt != nil && !strings.HasSuffix(u.Email, legacyPendingEmailSuffix)
}

func (s *Service) NeedsEmailSetup(u *domain.User) bool {
	if s.IsTelegramLinked(u) {
		return false
	}
	return strings.HasSuffix(u.Email, legacyPendingEmailSuffix) || u.EmailVerifiedAt == nil
}

func (s *Service) NeedsTelegramSetup(u *domain.User) bool {
	return !s.IsTelegramLinked(u)
}

func (s *Service) UserMeFromUser(u *domain.User) domain.UserMe {
	public := domain.UserPublic{
		ID:        u.ID.String(),
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
	}
	email := u.Email
	if strings.HasSuffix(email, legacyPendingEmailSuffix) {
		email = ""
	}
	tgUser := ""
	if u.TelegramUsername != nil {
		tgUser = *u.TelegramUsername
	}
	return domain.UserMe{
		UserPublic:         public,
		Email:              email,
		EmailVerified:      s.IsEmailVerified(u),
		NeedsEmailSetup:    s.NeedsEmailSetup(u),
		TelegramLinked:     s.IsTelegramLinked(u),
		TelegramUsername:   tgUser,
		NeedsTelegramSetup: s.NeedsTelegramSetup(u),
	}
}

func (s *Service) issueToken(u *domain.User) (string, error) {
	return s.jwt.GenerateToken(u.ID, u.Username, s.IsEmailVerified(u))
}

func generateSixDigitCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func hashCode(code string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	return string(b), err
}

func checkCode(hash, code string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(code)) == nil
}

func (s *Service) sendCodeEmail(to, code, subject string) error {
	if s.mail == nil || !s.mail.Enabled() {
		return ErrEmailNotConfigured
	}
	return s.mail.SendVerificationCode(to, code, subject)
}

func (s *Service) checkEmailRateLimit(ctx context.Context, email string) error {
	since := time.Now().Add(-emailRateLimitWindow)
	n, err := s.storage.CountAuthCodesSince(ctx, email, since)
	if err != nil {
		return err
	}
	if n >= emailRateLimitMax {
		return ErrEmailRateLimit
	}
	return nil
}

func (s *Service) createAndSendCode(ctx context.Context, email, purpose, subject string, userID *uuid.UUID, username, passwordHash *string) error {
	if err := s.checkEmailRateLimit(ctx, email); err != nil {
		return err
	}
	if err := s.storage.InvalidateAuthCodes(ctx, email, purpose); err != nil {
		return err
	}

	code, err := generateSixDigitCode()
	if err != nil {
		return err
	}
	codeHash, err := hashCode(code)
	if err != nil {
		return err
	}

	row := &domain.AuthCode{
		ID:           uuid.New(),
		Email:        email,
		CodeHash:     codeHash,
		Purpose:      purpose,
		UserID:       userID,
		Username:     username,
		PasswordHash: passwordHash,
		ExpiresAt:    time.Now().Add(authCodeTTL),
	}
	if err := s.storage.CreateAuthCode(ctx, row); err != nil {
		return err
	}
	// temporary disable send email verification code
	// render blocks stmp trafic 
	// if err := s.sendCodeEmail(email, code, subject); err != nil {
	// 	return err
	// }
	return nil
}

func (s *Service) StartSignup(ctx context.Context, email, username, password string) error {
	email, err := normalizeEmail(email)
	if err != nil {
		return err
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if password == "" {
		return fmt.Errorf("password is empty")
	}

	taken, err := s.storage.IsEmailTaken(ctx, email, nil)
	if err != nil {
		return err
	}
	if taken {
		return ErrEmailTaken
	}

	uid := uuid.Nil
	takenName, err := s.storage.IsUsernameTaken(ctx, username, uid)
	if err != nil {
		return err
	}
	if takenName {
		return ErrUsernameTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}
	ph := string(passwordHash)
	un := username

	return s.createAndSendCode(ctx, email, authPurposeSignup, "Регистрация — Eurovision Voting", nil, &un, &ph)
}

func (s *Service) ConfirmSignup(ctx context.Context, email, code string) (string, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return "", err
	}
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return "", ErrInvalidCode
	}

	ac, err := s.findAuthCode(ctx, email, authPurposeSignup, code)
	if err != nil {
		return "", err
	}
	if ac.Username == nil || ac.PasswordHash == nil {
		return "", ErrInvalidCode
	}

	now := time.Now()
	user := &domain.User{
		ID:              uuid.New(),
		Username:        *ac.Username,
		Email:           email,
		EmailVerifiedAt: &now,
		HashedPassword:  *ac.PasswordHash,
	}
	if err := s.storage.CreateUser(ctx, user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "email") || strings.Contains(pgErr.Message, "email") {
				return "", ErrEmailTaken
			}
			return "", ErrUsernameTaken
		}
		return "", err
	}
	_ = s.storage.MarkAuthCodeUsed(ctx, ac.ID)

	return s.issueToken(user)
}

func (s *Service) findAuthCode(ctx context.Context, email, purpose, code string) (*domain.AuthCode, error) {
	rows, err := s.storage.ListActiveAuthCodes(ctx, email, purpose)
	if err != nil {
		return nil, err
	}
	for _, ac := range rows {
		if checkCode(ac.CodeHash, code) {
			return ac, nil
		}
	}
	return nil, ErrInvalidCode
}

func (s *Service) SignIn(ctx context.Context, identifier, password string) (string, *domain.UserMe, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" || password == "" {
		return "", nil, fmt.Errorf("credentials required")
	}

	var u *domain.User
	var err error
	if strings.Contains(identifier, "@") {
		email, e := normalizeEmail(identifier)
		if e != nil {
			return "", nil, e
		}
		u, err = s.storage.GetUserByEmail(ctx, email)
	} else {
		u, err = s.storage.GetUserByUsername(ctx, identifier)
	}
	if err != nil {
		return "", nil, ErrUserNotExists
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)); err != nil {
		return "", nil, ErrWrongPassword
	}

	token, err := s.issueToken(u)
	if err != nil {
		return "", nil, err
	}
	me := s.UserMeFromUser(u)
	return token, &me, nil
}

func (s *Service) RequestEmailBind(ctx context.Context, userID uuid.UUID, email string) error {
	email, err := normalizeEmail(email)
	if err != nil {
		return err
	}

	if _, err := s.storage.GetUser(ctx, userID); err != nil {
		return err
	}

	taken, err := s.storage.IsEmailTaken(ctx, email, &userID)
	if err != nil {
		return err
	}
	if taken {
		return ErrEmailTaken
	}

	uid := userID
	return s.createAndSendCode(ctx, email, authPurposeEmailBind, "Подтверждение email — Eurovision Voting", &uid, nil, nil)
}

func (s *Service) ConfirmEmailBind(ctx context.Context, userID uuid.UUID, email, code string) (string, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return "", err
	}
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return "", ErrInvalidCode
	}

	ac, err := s.findAuthCode(ctx, email, authPurposeEmailBind, code)
	if err != nil {
		return "", err
	}
	if ac.UserID == nil || *ac.UserID != userID {
		return "", ErrInvalidCode
	}

	if err := s.storage.UpdateUserEmail(ctx, userID, email, true); err != nil {
		return "", err
	}
	_ = s.storage.MarkAuthCodeUsed(ctx, ac.ID)

	u, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}
	return s.issueToken(u)
}

func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	email, err := normalizeEmail(email)
	if err != nil {
		return err
	}

	u, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		// Do not reveal whether email exists
		return nil
	}
	if strings.HasSuffix(u.Email, legacyPendingEmailSuffix) {
		return nil
	}

	uid := u.ID
	return s.createAndSendCode(ctx, email, authPurposePasswordReset, "Сброс пароля — Eurovision Voting", &uid, nil, nil)
}

func (s *Service) ConfirmPasswordReset(ctx context.Context, email, code, newPassword string) error {
	email, err := normalizeEmail(email)
	if err != nil {
		return err
	}
	
	code = strings.TrimSpace(code)

	ac, err := s.findAuthCode(ctx, email, authPurposePasswordReset, code)
	if err != nil {
		return err
	}
	if ac.UserID == nil {
		return ErrInvalidCode
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.storage.UpdateUserPassword(ctx, *ac.UserID, string(passwordHash)); err != nil {
		return err
	}
	return s.storage.MarkAuthCodeUsed(ctx, ac.ID)
}

func (s *Service) ChangeUsername(ctx context.Context, userID uuid.UUID, username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 2 {
		return "", fmt.Errorf("username too short")
	}

	taken, err := s.storage.IsUsernameTaken(ctx, username, userID)
	if err != nil {
		return "", err
	}
	if taken {
		return "", ErrUsernameTaken
	}

	if err := s.storage.UpdateUsername(ctx, userID, username); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", ErrUsernameTaken
		}
		return "", err
	}

	u, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}
	return s.issueToken(u)
}
