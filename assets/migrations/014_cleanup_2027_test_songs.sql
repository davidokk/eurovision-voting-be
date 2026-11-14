-- Очистка тестовых песен конкурса 2027 из миграции 013_seed_2027_test_songs.sql.

DELETE FROM scores
WHERE performance_id IN (
    'b0000000-2027-4000-8000-000000000001'::uuid,
    'b0000000-2027-4000-8000-000000000002'::uuid
);

DELETE FROM messages
WHERE performance_id IN (
    'b0000000-2027-4000-8000-000000000001'::uuid,
    'b0000000-2027-4000-8000-000000000002'::uuid
);

DELETE FROM performance
WHERE id IN (
    'b0000000-2027-4000-8000-000000000001'::uuid,
    'b0000000-2027-4000-8000-000000000002'::uuid
);

-- Удаляем тестовый конкурс, если выступлений больше не осталось.
DELETE FROM contests
WHERE id = 'a0000000-2027-4000-8000-000000000001'::uuid
  AND NOT EXISTS (
      SELECT 1
      FROM performance
      WHERE contest_id = 'a0000000-2027-4000-8000-000000000001'::uuid
  );
