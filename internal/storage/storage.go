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
			(id, created_at, username, password)
		values
			($1, $2, $3, $4)
	`
	_, err := s.pool.Exec(ctx, query, user.ID, time.Now(), user.Username, user.HashedPassword)
	return err
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `select id, password, role from users where username = $1`
	row := s.pool.QueryRow(ctx, query, username)
	u := &domain.User{
		Username: username,
	}
	if err := row.Scan(&u.ID, &u.HashedPassword, &u.Role); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return u, nil
}

func (s *Storage) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `select id, username, password, role from users where id = $1`
	row := s.pool.QueryRow(ctx, query, id)
	u := &domain.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.HashedPassword, &u.Role); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return u, nil
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

func (s *Storage) RatePerformance(ctx context.Context, userID, performanceID uuid.UUID, score int, comment string, gif string) error {
	query := `
		insert into scores 
			(user_id, performance_id, score, comment, gif_url)
		values 
			($1, $2, $3, $4, $5)
		on conflict (user_id, performance_id) do update set
			score = $3,
			comment = CASE
				WHEN EXCLUDED.comment <> ''
                THEN EXCLUDED.comment
                ELSE scores.comment
			END,
			gif_url = CASE
				WHEN EXCLUDED.gif_url <> ''
                THEN EXCLUDED.gif_url
                ELSE scores.gif_url
    		END
	`
	_, err := s.pool.Exec(ctx, query, userID, performanceID, score, comment, gif)
	return err
}

func (s *Storage) UpdatePerformance(ctx context.Context, performanceID uuid.UUID, qualified bool, link string) error {
	_, err := s.pool.Exec(
		ctx, "update performance set qualified = $1, youtube_link = $2 where id = $3",
		qualified, link, performanceID,
	)
	return err
}

func (s *Storage) GetScoresFiltered(ctx context.Context, f domain.Filters) ([]domain.ScoreFiltered, error) {
	var (
		args  []any
		where []string
	)

	query := `
WITH base AS (
    SELECT 
        u.username,
        c.id as country_id,
        c.name_ru as country_name,
        co.year as contest_year,
        co.type as contest_type,
        sc.score,
        sc.comment,
        p.youtube_link as youtube_link,
        sc.gif_url,
        p.song as song,
        p.artist as artist,
        p.number as performance_number,
        sc.user_id
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

	query += `
),
ranked AS (
    SELECT *,
        ROW_NUMBER() OVER (
            PARTITION BY country_id, contest_year
            ORDER BY 
                CASE WHEN contest_type = 'final' THEN 1 ELSE 2 END,
                score DESC,
                performance_number ASC
        ) as rn
    FROM base
)
SELECT 
    username,
    country_name,
    contest_year,
    contest_type,
    score,
    comment,
    youtube_link,
    gif_url,
    song,
    artist
FROM ranked
WHERE rn = 1
`

	// -------------------------
	// SORTING (поверх уже очищенного набора)
	// -------------------------
	switch f.Sort {
	case domain.SortByScore:
		query += " ORDER BY score DESC, contest_year DESC, performance_number ASC"
	case domain.SortByTime:
		query += " ORDER BY contest_year DESC, performance_number ASC"
	default:
		query += " ORDER BY contest_year DESC, performance_number ASC"
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

func (s *Storage) InsertMessage(ctx context.Context, msg *domain.Message) error {
	query := `
		INSERT INTO messages (
			contest_id,
			user_id,
			message,
			created_at
		)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.pool.Exec(
		ctx,
		query,
		msg.ContestID,
		msg.UserID,
		msg.Message,
		msg.CreatedAt,
	)
	return err
}

func (s *Storage) GetMessages(ctx context.Context, contestID uuid.UUID) ([]domain.Message, error) {
	query := `
		SELECT 
			m.contest_id,
			m.user_id,
			m.message,
			m.created_at,
			u.username
		FROM messages m
			JOIN users u on u.id = m.user_id
		WHERE contest_id = $1
		ORDER BY created_at ASC
	`
	rows, err := s.pool.Query(ctx, query, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []domain.Message{}

	for rows.Next() {
		var m domain.Message

		if err := rows.Scan(
			&m.ContestID,
			&m.UserID,
			&m.Message,
			&m.CreatedAt,
			&m.Username,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}