-- Eurovision Song Contest 2021 — Grand Final (Rotterdam)
-- Running order, artists, songs and official ESC YouTube live performances.

INSERT INTO countries (id, name_ru, flag_emogi)
SELECT gen_random_uuid(), 'Россия', '🇷🇺'
WHERE NOT EXISTS (SELECT 1 FROM countries WHERE name_ru = 'Россия');

INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2021,
    '2021-05-22 21:00:00+02',
    '2021-05-23 01:00:00+02'
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
-- 1
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 1, 'https://www.youtube.com/watch?v=MgGc1F_cOE8', 'Elena Tsagrinou', 'El Diablo'),
-- 2
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Албания'), 2, 'https://www.youtube.com/watch?v=OxniMKVXjqM', 'Anxhela Peristeri', 'Karma'),
-- 3
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Израиль'), 3, 'https://www.youtube.com/watch?v=5Y-blG5ptwo', 'Eden Alene', 'Set Me Free'),
-- 4
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Бельгия'), 4, 'https://www.youtube.com/watch?v=HbpxcUMtjwY', 'Hooverphonic', 'The Wrong Place'),
-- 5
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Россия'), 5, 'https://www.youtube.com/watch?v=V-Di9A28e5E', 'Manizha', 'Russian Woman'),
-- 6
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Мальта'), 6, 'https://www.youtube.com/watch?v=jbbc2yKnv0M', 'Destiny', 'Je Me Casse'),
-- 7
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Португалия'), 7, 'https://www.youtube.com/watch?v=qq9SceXaukA', 'The Black Mamba', 'Love Is On My Side'),
-- 8
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 8, 'https://www.youtube.com/watch?v=Q1lkpIqf7a4', 'Hurricane', 'Loco Loco'),
-- 9
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 9, 'https://www.youtube.com/watch?v=VxNOynEJ6wc', 'James Newman', 'Embers'),
-- 10
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Греция'), 10, 'https://www.youtube.com/watch?v=QhZhJFp0vVw', 'Stefania', 'Last Dance'),
-- 11
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Швейцария'), 11, 'https://www.youtube.com/watch?v=jznH_fltcYA', 'Gjon''s Tears', 'Tout l''Univers'),
-- 12
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Исландия'), 12, 'https://www.youtube.com/watch?v=YSMhu-PrLME', 'Daði og Gagnamagnið', '10 Years'),
-- 13
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Испания'), 13, 'https://www.youtube.com/watch?v=thV-vhOQHvQ', 'Blas Cantó', 'Voy A Quedarme'),
-- 14
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Молдова'), 14, 'https://www.youtube.com/watch?v=5I_1GTehgkI', 'Natalia Gordienko', 'SUGAR'),
-- 15
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Германия'), 15, 'https://www.youtube.com/watch?v=1m0VEAfLV4E', 'Jendrik', 'I Don''t Feel Hate'),
-- 16
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Финляндия'), 16, 'https://www.youtube.com/watch?v=fOPGxQ4fgVw', 'Blind Channel', 'Dark Side'),
-- 17
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Болгария'), 17, 'https://www.youtube.com/watch?v=B1ZiPrMwMYA', 'VICTORIA', 'Growing Up Is Getting Old'),
-- 18
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Литва'), 18, 'https://www.youtube.com/watch?v=JNweec5olYw', 'The Roop', 'Discoteque'),
-- 19
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Украина'), 19, 'https://www.youtube.com/watch?v=lqvzDkgok_g', 'Go_A', 'Shum'),
-- 20
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Франция'), 20, 'https://www.youtube.com/watch?v=Unj9WbeLzRU', 'Barbara Pravi', 'Voilà'),
-- 21
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 21, 'https://www.youtube.com/watch?v=FTQ22S5YC7Q', 'Efendi', 'Mata Hari'),
-- 22
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 22, 'https://www.youtube.com/watch?v=F0r3vmGLzZU', 'TIX', 'Fallen Angel'),
-- 23
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Нидерланды'), 23, 'https://www.youtube.com/watch?v=o3Da5cr2fTs', 'Jeangu Macrooy', 'Birth Of A New Age'),
-- 24
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Италия'), 24, 'https://www.youtube.com/watch?v=RVH5dn1cxAQ', 'Måneskin', 'Zitti E Buoni'),
-- 25
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 25, 'https://www.youtube.com/watch?v=NyJRyif6_kk', 'Tusse', 'Voices'),
-- 26
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2021), (SELECT id FROM countries WHERE name_ru = 'Сан-Марино'), 26, 'https://www.youtube.com/watch?v=TkhiH-JXFPs', 'Senhit', 'Adrenalina');
