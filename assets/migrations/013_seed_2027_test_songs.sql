-- Тестовые данные для конкурса 2027 (проверка нового функционала).
-- Удаляются миграцией 014_cleanup_2027_test_songs.sql.

INSERT INTO contests (id, type, year, starts, ends)
SELECT
    'a0000000-2027-4000-8000-000000000001'::uuid,
    'first-semifinal',
    2027,
    now() - interval '1 hour',
    now() + interval '30 days'
WHERE NOT EXISTS (
    SELECT 1 FROM contests WHERE id = 'a0000000-2027-4000-8000-000000000001'::uuid
);

INSERT INTO performance (
    id,
    contest_id,
    country_id,
    number,
    youtube_link,
    artist,
    song
)
VALUES
(
    'b0000000-2027-4000-8000-000000000001'::uuid,
    'a0000000-2027-4000-8000-000000000001'::uuid,
    (SELECT id FROM countries WHERE name_ru = 'Швеция' LIMIT 1),
    1,
    'https://www.youtube.com/watch?v=dQw4w9WgXcQ',
    'Test Artist One',
    'Test Song Alpha 2027'
),
(
    'b0000000-2027-4000-8000-000000000002'::uuid,
    'a0000000-2027-4000-8000-000000000001'::uuid,
    (SELECT id FROM countries WHERE name_ru = 'Украина' LIMIT 1),
    2,
    'https://www.youtube.com/watch?v=9bZkp7q19f0',
    'Test Artist Two',
    'Test Song Beta 2027'
)
ON CONFLICT (id) DO NOTHING;
