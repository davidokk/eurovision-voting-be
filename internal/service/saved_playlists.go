package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/storage"

	"github.com/google/uuid"
)

const maxSavedPlaylistEntries = 100
const maxSavedPlaylistNameLen = 80

func normalizePlaylistName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("name required")
	}
	if utf8.RuneCountInString(name) > maxSavedPlaylistNameLen {
		return "", fmt.Errorf("name too long")
	}
	return name, nil
}

func validateSavedPlaylistEntries(entries []domain.GamePlaylistEntryInput) error {
	if len(entries) == 0 {
		return fmt.Errorf("playlist must have at least one track")
	}
	if len(entries) > maxSavedPlaylistEntries {
		return fmt.Errorf("too many tracks (max %d)", maxSavedPlaylistEntries)
	}
	for _, e := range entries {
		id := strings.TrimSpace(e.PerformanceID)
		link := strings.TrimSpace(e.YoutubeLink)
		if id == "" && link == "" {
			return fmt.Errorf("invalid playlist entry")
		}
	}
	return nil
}

func (s *Service) ListSavedPlaylists(ctx context.Context, userID uuid.UUID) ([]domain.SavedPlaylistSummary, error) {
	return s.storage.ListSavedPlaylists(ctx, userID)
}

func (s *Service) GetSavedPlaylist(ctx context.Context, userID, playlistID uuid.UUID) (*domain.SavedPlaylist, error) {
	item, err := s.storage.GetSavedPlaylist(ctx, userID, playlistID)
	if errors.Is(err, storage.ErrSavedPlaylistNotFound) {
		return nil, ErrGameNotFound
	}
	return item, err
}

func (s *Service) CreateSavedPlaylist(ctx context.Context, userID uuid.UUID, name string, entries []domain.GamePlaylistEntryInput) (*domain.SavedPlaylist, error) {
	name, err := normalizePlaylistName(name)
	if err != nil {
		return nil, err
	}
	if err := validateSavedPlaylistEntries(entries); err != nil {
		return nil, err
	}
	return s.storage.CreateSavedPlaylist(ctx, userID, name, entries)
}

func (s *Service) UpdateSavedPlaylist(ctx context.Context, userID, playlistID uuid.UUID, name string, entries []domain.GamePlaylistEntryInput) (*domain.SavedPlaylist, error) {
	name, err := normalizePlaylistName(name)
	if err != nil {
		return nil, err
	}
	if err := validateSavedPlaylistEntries(entries); err != nil {
		return nil, err
	}
	item, err := s.storage.UpdateSavedPlaylist(ctx, userID, playlistID, name, entries)
	if errors.Is(err, storage.ErrSavedPlaylistNotFound) {
		return nil, ErrGameNotFound
	}
	return item, err
}

func (s *Service) DeleteSavedPlaylist(ctx context.Context, userID, playlistID uuid.UUID) error {
	err := s.storage.DeleteSavedPlaylist(ctx, userID, playlistID)
	if errors.Is(err, storage.ErrSavedPlaylistNotFound) {
		return ErrGameNotFound
	}
	return err
}
