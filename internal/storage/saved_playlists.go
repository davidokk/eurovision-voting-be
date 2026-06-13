package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrSavedPlaylistNotFound = errors.New("saved playlist not found")

func (s *Storage) ListSavedPlaylists(ctx context.Context, userID uuid.UUID) ([]domain.SavedPlaylistSummary, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, jsonb_array_length(entries), updated_at
		FROM saved_playlists
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.SavedPlaylistSummary
	for rows.Next() {
		var item domain.SavedPlaylistSummary
		if err := rows.Scan(&item.ID, &item.Name, &item.EntryCount, &item.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if out == nil {
		out = []domain.SavedPlaylistSummary{}
	}
	return out, rows.Err()
}

func (s *Storage) GetSavedPlaylist(ctx context.Context, userID, playlistID uuid.UUID) (*domain.SavedPlaylist, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, name, entries, created_at, updated_at
		FROM saved_playlists
		WHERE id = $1 AND user_id = $2
	`, playlistID, userID)

	var item domain.SavedPlaylist
	var raw []byte
	if err := row.Scan(&item.ID, &item.Name, &raw, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSavedPlaylistNotFound
		}
		return nil, err
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &item.Entries); err != nil {
			return nil, err
		}
	}
	if item.Entries == nil {
		item.Entries = []domain.GamePlaylistEntryInput{}
	}
	return &item, nil
}

func (s *Storage) CreateSavedPlaylist(ctx context.Context, userID uuid.UUID, name string, entries []domain.GamePlaylistEntryInput) (*domain.SavedPlaylist, error) {
	raw, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}

	var item domain.SavedPlaylist
	var id uuid.UUID
	err = s.pool.QueryRow(ctx, `
		INSERT INTO saved_playlists (user_id, name, entries)
		VALUES ($1, $2, $3::jsonb)
		RETURNING id, name, entries, created_at, updated_at
	`, userID, name, raw).Scan(&id, &item.Name, &raw, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "saved_playlists_user_name_idx") {
			return nil, fmt.Errorf("playlist name already exists")
		}
		return nil, err
	}
	item.ID = id.String()
	if err := json.Unmarshal(raw, &item.Entries); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Storage) UpdateSavedPlaylist(ctx context.Context, userID, playlistID uuid.UUID, name string, entries []domain.GamePlaylistEntryInput) (*domain.SavedPlaylist, error) {
	raw, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}

	var item domain.SavedPlaylist
	var rawOut []byte
	err = s.pool.QueryRow(ctx, `
		UPDATE saved_playlists
		SET name = $3, entries = $4::jsonb, updated_at = $5
		WHERE id = $1 AND user_id = $2
		RETURNING id, name, entries, created_at, updated_at
	`, playlistID, userID, name, raw, time.Now().UTC()).Scan(&item.ID, &item.Name, &rawOut, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSavedPlaylistNotFound
		}
		if strings.Contains(err.Error(), "saved_playlists_user_name_idx") {
			return nil, fmt.Errorf("playlist name already exists")
		}
		return nil, err
	}
	if err := json.Unmarshal(rawOut, &item.Entries); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Storage) DeleteSavedPlaylist(ctx context.Context, userID, playlistID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM saved_playlists WHERE id = $1 AND user_id = $2
	`, playlistID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSavedPlaylistNotFound
	}
	return nil
}
