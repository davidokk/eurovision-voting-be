package storage

import (
	"context"
	"encoding/json"
	"eurovision-voting/internal/domain"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg Config) (*Storage, error) {
	databaseURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	databaseURL.User = url.UserPassword(cfg.Username, cfg.Password)

	connConfig, err := pgxpool.ParseConfig(databaseURL.String())
	if err != nil {
		return nil, err
	}
	connConfig.MinConns = cfg.MinConns
	connConfig.MaxConns = cfg.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		insert into users 
			(id, created_at, username, password, email, email_verified_at)
		values
			($1, $2, $3, $4, lower(trim($5)), $6)
	`
	_, err := s.pool.Exec(ctx, query, user.ID, time.Now(), user.Username, user.HashedPassword, user.Email, user.EmailVerifiedAt)
	return err
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `select id, username, password, role, avatar_url, coalesce(email, ''), email_verified_at from users where username = $1`
	row := s.pool.QueryRow(ctx, query, username)
	return scanUserRow(row, username)
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `select id, username, password, role, avatar_url, coalesce(email, ''), email_verified_at from users where lower(trim(email)) = lower(trim($1))`
	row := s.pool.QueryRow(ctx, query, email)
	return scanUserRow(row, "")
}

func (s *Storage) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `select id, username, password, role, avatar_url, coalesce(email, ''), email_verified_at from users where id = $1`
	row := s.pool.QueryRow(ctx, query, id)
	return scanUserRow(row, "")
}

func scanUserRow(row interface {
	Scan(dest ...any) error
}, usernameFallback string) (*domain.User, error) {
	u := &domain.User{}
	if usernameFallback != "" {
		u.Username = usernameFallback
	}
	if err := row.Scan(&u.ID, &u.Username, &u.HashedPassword, &u.Role, &u.AvatarURL, &u.Email, &u.EmailVerifiedAt); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return u, nil
}

func (s *Storage) UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string, verified bool) error {
	if verified {
		_, err := s.pool.Exec(ctx, `
			UPDATE users SET email = lower(trim($2)), email_verified_at = now() WHERE id = $1
		`, userID, email)
		return err
	}
	_, err := s.pool.Exec(ctx, `
		UPDATE users SET email = lower(trim($2)), email_verified_at = NULL WHERE id = $1
	`, userID, email)
	return err
}

func (s *Storage) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	_, err := s.pool.Exec(ctx, `UPDATE users SET password = $2 WHERE id = $1`, userID, passwordHash)
	return err
}

func (s *Storage) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	_, err := s.pool.Exec(ctx, `UPDATE users SET username = $2 WHERE id = $1`, userID, username)
	return err
}

func (s *Storage) IsUsernameTaken(ctx context.Context, username string, excludeUserID uuid.UUID) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id <> $2)
	`, username, excludeUserID).Scan(&exists)
	return exists, err
}

func (s *Storage) IsEmailTaken(ctx context.Context, email string, excludeUserID *uuid.UUID) (bool, error) {
	var exists bool
	if excludeUserID != nil {
		err := s.pool.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM users
				WHERE lower(trim(email)) = lower(trim($1))
				  AND id <> $2
				  AND email NOT LIKE '%@legacy.pending'
			)
		`, email, *excludeUserID).Scan(&exists)
		return exists, err
	}
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE lower(trim(email)) = lower(trim($1))
			  AND email NOT LIKE '%@legacy.pending'
		)
	`, email).Scan(&exists)
	return exists, err
}

func (s *Storage) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error {
	_, err := s.pool.Exec(ctx, `UPDATE users SET avatar_url = $2 WHERE id = $1`, userID, avatarURL)
	return err
}

func (s *Storage) ClearUserAvatar(ctx context.Context, userID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `UPDATE users SET avatar_url = NULL WHERE id = $1`, userID)
	return err
}

func (s *Storage) GetContestList(ctx context.Context) ([]domain.Contest, error) {
	rows, err := s.pool.Query(ctx, "select id, year, type from contests order by type")
	if err != nil {
		return nil, err
	}
	res := []domain.Contest{}
	for rows.Next() {
		c := domain.Contest{}
		if err := rows.Scan(&c.ID, &c.Year, &c.Type); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		res = append(res, c)
	}
	return res, nil
}

func (s *Storage) GetContestByPerformance(ctx context.Context, id uuid.UUID) (*domain.Contest, error) {
	query := `
		SELECT c.id, c.type, c.year, c.starts, c.ends
		FROM contests c
		JOIN performance p ON p.contest_id = c.id
		WHERE p.id = $1
		LIMIT 1
	`
	var c domain.Contest
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Type,
		&c.Year,
		&c.Starts,
		&c.Ends,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Storage) GetContestView(ctx context.Context, contestID uuid.UUID) (*domain.ContestView, error) {
	query := `
		SELECT
			json_build_object(
				'id', c.id,
				'type', c.type,
				'year', c.year,
				'starts', c.starts,
				'ends', c.ends
			) AS contest,

			COALESCE(
				json_agg(
					json_build_object(
						'performance_id', p.id,
						'qualified', p.qualified,

						'country', json_build_object(
							'id', co.id,
							'name_ru', co.name_ru,
							'flag_emoji', co.flag_emogi
						),

						'artist', p.artist,
						'song', p.song,
						'number', p.number,
						'youtube_link', p.youtube_link,
						'place', p.place,

						'total_score', COALESCE(s.total_score, 0),

						'scores', COALESCE(sc.scores, '[]'::json)
					)
					ORDER BY p.number
				) FILTER (WHERE p.id IS NOT NULL),
				'[]'::json
			) AS performances

		FROM contests c

		LEFT JOIN performance p
			ON p.contest_id = c.id

		LEFT JOIN countries co
			ON co.id = p.country_id

		LEFT JOIN (
			SELECT
				performance_id,
				AVG(score) AS total_score
			FROM scores
			GROUP BY performance_id
		) s ON s.performance_id = p.id

		LEFT JOIN (
			SELECT
				performance_id,
				json_agg(
					json_build_object(
						'user_id', s.user_id,
						'username', u.username,
						'score', s.score,
						'comment', s.comment,
						'gif_url', s.gif_url
					)
				) AS scores
			FROM scores s JOIN users u on s.user_id = u.id
			GROUP BY performance_id
		) sc ON sc.performance_id = p.id

		WHERE c.id = $1

		GROUP BY c.id, c.type, c.year;
	`

	var result struct {
		Contest      json.RawMessage `json:"contest"`
		Performances json.RawMessage `json:"performances"`
	}
	if err := s.pool.QueryRow(ctx, query, contestID).Scan(
		&result.Contest,
		&result.Performances,
	); err != nil {
		return nil, err
	}

	var contest domain.Contest
	if err := json.Unmarshal(result.Contest, &contest); err != nil {
		return nil, err
	}
	var performances []domain.PerformanceWithScores
	if err := json.Unmarshal(result.Performances, &performances); err != nil {
		return nil, err
	}

	return &domain.ContestView{
		Contest:      contest,
		Performances: performances,
	}, nil
}

func (s *Storage) RatePerformance(ctx context.Context, userID, performanceID uuid.UUID, score int, comment string, gif string) (int, error) {
	query := `
		with old as (
			select score as old_score
			from scores
			where user_id = $1 and performance_id = $2
		),
		upsert as (
			insert into scores (user_id, performance_id, score, comment, gif_url)
			values ($1, $2, $3, $4, $5)
			on conflict (user_id, performance_id) do update set
				score = $3,
				comment = case
					when excluded.comment <> '' then excluded.comment
					else scores.comment
				end,
				gif_url = case
					when excluded.gif_url <> '' then excluded.gif_url
					else scores.gif_url
				end
			returning score
		)
		select old.old_score
		from old;
	`
	row := s.pool.QueryRow(ctx, query, userID, performanceID, score, comment, gif)
	var old int 
	if err := row.Scan(&old); err != nil {
		return 0, nil
	}
	return old, nil
}

func (s *Storage) UpdatePerformance(ctx context.Context, performanceID uuid.UUID, qualified bool, link string, place int) error {
	_, err := s.pool.Exec(ctx, `
		update performance 
		set 
			qualified = $1, 
			youtube_link = $2,
			place = $3 
		where id = $4
		`, qualified, link, place, performanceID,
	)
	return err
}

func (s *Storage) GetScoresFiltered(ctx context.Context, f domain.Filters) ([]domain.ScoreFiltered, error) {
	var (
		args  []any
		where []string
	)

	query := `
    SELECT 
        u.username,
        c.name_ru as country_name,
        co.year as contest_year,
        co.type as contest_type,
        sc.score,
        sc.comment,
        p.youtube_link as youtube_link,
        sc.gif_url,
        p.song as song,
        p.artist as artist,
		p.qualified as qualified,
		p.place
    FROM scores sc
    JOIN users u ON u.id = sc.user_id
    JOIN performance p ON p.id = sc.performance_id
    JOIN countries c ON c.id = p.country_id
    JOIN contests co ON co.id = p.contest_id
    WHERE 1=1
`
	// -------------------------
	// FILTERS
	// -------------------------
	if f.UserID != nil {
		args = append(args, *f.UserID)
		where = append(where, fmt.Sprintf("sc.user_id = $%d", len(args)))
	}

	if f.CountryID != nil {
		args = append(args, *f.CountryID)
		where = append(where, fmt.Sprintf("p.country_id = $%d", len(args)))
	}

	if f.ContestYear != nil {
		args = append(args, *f.ContestYear)
		where = append(where, fmt.Sprintf("co.year = $%d", len(args)))
	}

	if len(where) > 0 {
		query += " AND " + strings.Join(where, " AND ")
	}

	// -------------------------
	// SORTING (поверх уже очищенного набора)
	// -------------------------
	switch f.Sort {
	case domain.SortByScore:
		query += " ORDER BY sc.score DESC, co.year DESC, p.number ASC"
	case domain.SortByTime:
		query += " ORDER BY co.year DESC, p.number ASC"
	default:
		query += " ORDER BY co.year DESC, p.number ASC"
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ScoreFiltered

	for rows.Next() {
		var r domain.ScoreFiltered

		err := rows.Scan(
			&r.Username,
			&r.CountryName,
			&r.ContestYear,
			&r.ContestType,
			&r.Score,
			&r.Comment,
			&r.YoutubeLink,
			&r.GifURL,
			&r.Song,
			&r.Artist,
			&r.Qualified,
			&r.Place,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, rows.Err()
}

func (s *Storage) GetCountries(ctx context.Context) ([]domain.Country, error) {
	rows, err := s.pool.Query(ctx, "select id, name_ru, flag_emogi from countries")
	if err != nil {
		return nil, err
	}
	res := []domain.Country{}
	for rows.Next() {
		var c domain.Country
		if err := rows.Scan(&c.ID, &c.NameRU, &c.FlagEmoji); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (s *Storage) GetCountryByPerformance(ctx context.Context, performanceID uuid.UUID) (*domain.Country, error) {
	rows := s.pool.QueryRow(ctx, `
		select c.id, c.name_ru, c.flag_emogi 
		from countries c
			join performance p on p.country_id = c.id
		where p.id = $1
		`,
		performanceID,
	)

	var c domain.Country
	if err := rows.Scan(&c.ID, &c.NameRU, &c.FlagEmoji); err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *Storage) InsertMessage(ctx context.Context, msg *domain.Message) error {
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	query := `
		INSERT INTO messages (
			id, type, performance_id, contest_id, user_id, message, created_at,
			score, old_score, comment, gif,
			reply_to_id, content_type, media_url, media_duration_ms
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`
	return s.pool.QueryRow(
		ctx,
		query,
		msg.ID,
		msg.Type,
		nullUUID(msg.PerformanceID),
		msg.ContestID,
		msg.UserID,
		msg.Message,
		msg.CreatedAt,
		msg.Score,
		msg.OldScore,
		msg.Comment,
		msg.Gif,
		msg.ReplyToID,
		msg.ContentType,
		msg.MediaURL,
		msg.MediaDurationMs,
	).Scan(&msg.ID)
}

func nullUUID(id uuid.UUID) any {
	if id == uuid.Nil {
		return nil
	}
	return id
}

// FillReplyPreview подгружает превью родительского сообщения для WS/ответа API.
func (s *Storage) FillReplyPreview(ctx context.Context, msg *domain.Message) error {
	if msg.ReplyToID == nil || *msg.ReplyToID == uuid.Nil {
		return nil
	}
	var parentUser, parentMsg, parentCT string
	var parentMedia *string
	err := s.pool.QueryRow(ctx, `
		SELECT pu.username, COALESCE(pm.message, ''), COALESCE(pm.content_type, 'text'), pm.media_url
		FROM messages pm
		JOIN users pu ON pu.id = pm.user_id
		WHERE pm.id = $1
	`, *msg.ReplyToID).Scan(&parentUser, &parentMsg, &parentCT, &parentMedia)
	if err != nil {
		return err
	}
	msg.ReplyTo = &domain.MessageReplyPreview{
		ID:          msg.ReplyToID.String(),
		Username:    parentUser,
		Message:     parentMsg,
		ContentType: parentCT,
		MediaURL:    parentMedia,
	}
	return nil
}

func (s *Storage) GetMessages(ctx context.Context, contestID uuid.UUID) ([]domain.Message, error) {
	query := `
		SELECT 
			m.id,
			COALESCE(m.type, ''),
			m.contest_id,
			m.user_id,
			m.message,
			m.created_at,
			u.username,
			u.avatar_url,
			c.name_ru,
			c.flag_emogi,
			m.score,
			m.old_score,
			m.comment,
			m.gif,
			COALESCE(m.content_type, 'text'),
			m.media_url,
			m.media_duration_ms,
			m.reply_to_id,
			pu.username,
			pm.message,
			COALESCE(pm.content_type, 'text'),
			pm.media_url
		FROM messages m
			JOIN users u ON u.id = m.user_id
			LEFT JOIN performance p ON p.id = m.performance_id 
			LEFT JOIN countries c ON c.id = p.country_id
			LEFT JOIN messages pm ON pm.id = m.reply_to_id
			LEFT JOIN users pu ON pu.id = pm.user_id
 		WHERE m.contest_id = $1
		ORDER BY m.created_at ASC
	`
	rows, err := s.pool.Query(ctx, query, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		var replyID *uuid.UUID
		var parentUser, parentMsg, parentCT *string
		var parentMedia *string

		if err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.ContestID,
			&m.UserID,
			&m.Message,
			&m.CreatedAt,
			&m.Username,
			&m.AvatarURL,
			&m.Country,
			&m.CountryFlag,
			&m.Score,
			&m.OldScore,
			&m.Comment,
			&m.Gif,
			&m.ContentType,
			&m.MediaURL,
			&m.MediaDurationMs,
			&replyID,
			&parentUser,
			&parentMsg,
			&parentCT,
			&parentMedia,
		); err != nil {
			return nil, err
		}
		m.ReplyToID = replyID
		if replyID != nil && parentUser != nil {
			m.ReplyTo = &domain.MessageReplyPreview{
				ID:          replyID.String(),
				Username:    *parentUser,
				Message:     strVal(parentMsg),
				ContentType: strVal(parentCT),
				MediaURL:    parentMedia,
			}
		}
		if m.ContentType == "" {
			m.ContentType = "text"
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func strVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
