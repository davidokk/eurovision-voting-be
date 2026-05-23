package service

import (
	"context"
	"fmt"
	"io"

	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/media"

	"github.com/google/uuid"
)

var (
	ErrMediaNotConfigured = fmt.Errorf("media storage not configured")
	ErrInvalidMediaKind   = fmt.Errorf("invalid media kind")
)

func (s *Service) UploadMedia(ctx context.Context, userID uuid.UUID, kind string, r io.Reader, contentType string) (string, error) {
	if s.s3 == nil {
		return "", ErrMediaNotConfigured
	}
	folder := fmt.Sprintf("%s/%s", kind, userID.String())
	maxW := media.MaxImageWidth
	if kind == "avatar" {
		maxW = media.MaxAvatarWidth
	}
	switch kind {
	case "avatar", "image":
		data, ct, err := media.CompressImage(r, maxW)
		if err != nil {
			return "", err
		}
		url, err := s.s3.Upload(ctx, folder, data, ct)
		if err != nil {
			return "", err
		}
		if kind == "avatar" {
			if err := s.storage.UpdateUserAvatar(ctx, userID, url); err != nil {
				return "", err
			}
		}
		return url, nil
	case "voice":
		data, ct, err := media.ValidateVoice(r)
		if err != nil {
			return "", err
		}
		if contentType != "" {
			ct = contentType
		}
		return s.s3.Upload(ctx, folder, data, ct)
	case "video_note":
		data, ct, err := media.ValidateVideoNote(r)
		if err != nil {
			return "", err
		}
		if contentType != "" {
			ct = contentType
		}
		return s.s3.Upload(ctx, folder, data, ct)
	default:
		return "", ErrInvalidMediaKind
	}
}

func (s *Service) GetUserPublic(ctx context.Context, id uuid.UUID) (*domain.UserPublic, error) {
	u, err := s.storage.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.UserPublic{
		ID:        u.ID.String(),
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
	}, nil
}

func (s *Service) GetUserMe(ctx context.Context, id uuid.UUID) (*domain.UserMe, error) {
	u, err := s.storage.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	me := s.UserMeFromUser(u)
	return &me, nil
}

func (s *Service) SetUserAvatar(ctx context.Context, userID uuid.UUID, url string) error {
	return s.storage.UpdateUserAvatar(ctx, userID, url)
}

func (s *Service) DeleteUserAvatar(ctx context.Context, userID uuid.UUID) error {
	return s.storage.ClearUserAvatar(ctx, userID)
}
