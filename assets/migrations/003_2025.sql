ALTER TABLE performance
ALTER COLUMN youtube_link DROP NOT NULL;

ALTER TABLE performance
ALTER COLUMN youtube_link SET DEFAULT '';

INSERT INTO contests (id, type, year, starts, ends)
VALUES
(
    gen_random_uuid(),
    'first-semifinal',
    2025,
    '2025-05-13 21:00:00+02',
    '2025-05-13 23:30:00+02'
),
(
    gen_random_uuid(),
    'second-semifinal',
    2025,
    '2025-05-15 21:00:00+02',
    '2025-05-15 23:30:00+02'
),
(
    gen_random_uuid(),
    'final',
    2025,
    '2025-05-17 21:00:00+02',
    '2025-05-18 01:00:00+02'
);

INSERT INTO countries (id, name_ru, flag_emogi) VALUES
(gen_random_uuid(), 'Исландия', '🇮🇸'),
(gen_random_uuid(), 'Словения', '🇸🇮'),
(gen_random_uuid(), 'Нидерланды', '🇳🇱'),
(gen_random_uuid(), 'Ирландия', '🇮🇪'),
(gen_random_uuid(), 'Испания', '🇪🇸');

INSERT INTO performance (
    id,
    contest_id,
    country_id,
    number,
    artist,
    song
)
VALUES

-- FIRST SEMIFINAL
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Исландия'), 1, 'Væb', 'Róa'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Польша'), 2, 'Justyna Steczkowska', 'Gaja'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Словения'), 3, 'Klemen', 'How Much Time Do We Have Left'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Эстония'), 4, 'Tommy Cash', 'Espresso macchiato'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Украина'), 5, 'Ziferblat', 'Bird of Pray'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Швеция'), 6, 'KAJ', 'Bara Bada Bastu'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Португалия'), 7, 'Napa', 'Deslocado'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Норвегия'), 8, 'Kyle Alessandro', 'Lighter'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Бельгия'), 9, 'Red Sebastian', 'Strobe Lights'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Азербайджан'), 10, 'Mamagama', 'Run With U'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Сан-Марино'), 11, 'Gabry Ponte', 'Tutta l''Italia'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Албания'), 12, 'Shkodra Elektronike', 'Zjerm'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Нидерланды'), 13, 'Claude', 'C''est La Vie'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Хорватия'), 14, 'Marko Bošnjak', 'Poison Cake'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='first-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Кипр'), 15, 'Theo Evan', 'Shh'),

-- SECOND SEMIFINAL
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Австралия'), 1, 'Go-Jo', 'Milkshake Man'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Черногория'), 2, 'Нина Жижич', 'Dobrodošli'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Ирландия'), 3, 'EMMY', 'Laika Party'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Латвия'), 4, 'Tautumeitas', 'Bur Man Laimi'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Армения'), 5, 'PARG', 'SURVIVOR'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Австрия'), 6, 'JJ', 'Wasted Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Греция'), 7, 'Klavdia', 'Asteromáta'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Литва'), 8, 'Katarsis', 'Tavo Akys'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Мальта'), 9, 'Мириана Конте', 'SERVING'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Грузия'), 10, 'Мариам Шенгелия', 'Свобода'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Дания'), 11, 'Сиссал', 'Галлюцинация'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Чехия'), 12, 'ADONXS', 'Kiss Kiss Goodbye'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Люксембург'), 13, 'Laura Thorn', 'La Poupée Monte Le Son'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Израиль'), 14, 'Yuval Raphael', 'New Day Will Rise'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Сербия'), 15, 'Princ', 'Mila'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='second-semifinal' AND year=2025), (SELECT id FROM countries WHERE name_ru='Финляндия'), 16, 'Erika Vikman', 'ICH KOMME'),

-- FINAL
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Норвегия'), 1, 'Kyle Alessandro', 'Lighter'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Люксембург'), 2, 'Laura Thorn', 'La Poupée Monte Le Son'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Эстония'), 3, 'Tommy Cash', 'Espresso macchiato'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Израиль'), 4, 'Yuval Raphael', 'New Day Will Rise'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Литва'), 5, 'Katarsis', 'Tavo Akys'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Испания'), 6, 'Melody', 'Esa Diva'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Украина'), 7, 'Ziferblat', 'Bird of Pray'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Великобритания'), 8, 'Remember Monday', 'What The Hell Just Happened?'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Австрия'), 9, 'JJ', 'Wasted Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Исландия'), 10, 'Væb', 'Róa'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Латвия'), 11, 'Tautumeitas', 'Bur Man Laimi'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Нидерланды'), 12, 'Claude', 'C''est La Vie'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Финляндия'), 13, 'Erika Vikman', 'ICH KOMME'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Италия'), 14, 'Lucio Corsi', 'Volevo Essere Un Duro'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Польша'), 15, 'Justyna Steczkowska', 'Gaja'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Германия'), 16, 'Abor & Tynna', 'Baller'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Греция'), 17, 'Klavdia', 'Asteromáta'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Армения'), 18, 'PARG', 'SURVIVOR'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Швейцария'), 19, 'Zoë Më', 'Voyage'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Мальта'), 20, 'Мириана Конте', 'SERVING'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Португалия'), 21, 'Napa', 'Deslocado'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Дания'), 22, 'Сиссал', 'Галлюцинация'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Швеция'), 23, 'KAJ', 'Bara Bada Bastu'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Франция'), 24, 'Луан', 'Маман'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Сан-Марино'), 25, 'Gabry Ponte', 'Tutta l''Italia'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type='final' AND year=2025), (SELECT id FROM countries WHERE name_ru='Албания'), 26, 'Shkodra Elektronike', 'Zjerm');

INSERT INTO scores (user_id, performance_id, score)
SELECT
    '85cfc6c3-5549-45cb-990f-86b29e0826c9',
    p.id,
    v.score
FROM (
    VALUES
        ('Исландия', 'Róa', 5),
        ('Польша', 'Gaja', 6),
        ('Словения', 'How Much Time Do We Have Left', 3),
        ('Эстония', 'Espresso macchiato', 10),
        ('Испания', 'Esa Diva', 9),
        ('Украина', 'Bird of Pray', 5),
        ('Швеция', 'Bara Bada Bastu', 9),
        ('Португалия', 'Deslocado', 2),
        ('Норвегия', 'Lighter', 10),
        ('Бельгия', 'Strobe Lights', 8),
        ('Италия', 'Volevo Essere Un Duro', 7),
        ('Азербайджан', 'Run With U', 4),
        ('Сан-Марино', 'Tutta l''Italia', 1),
        ('Албания', 'Zjerm', 7),
        ('Нидерланды', 'C''est La Vie', 5),
        ('Хорватия', 'Poison Cake', 8),
        ('Швейцария', 'Voyage', 8),
        ('Кипр', 'Shh', 10)
) AS v(country_name, song, score)
JOIN countries c
    ON c.name_ru = v.country_name
JOIN performance p
    ON p.country_id = c.id
   AND p.song = v.song
JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'final'
   AND ct.year = 2025;


INSERT INTO scores (user_id, performance_id, score)
SELECT
    '7c6bc8f1-1d53-4840-af70-f904a52745cc',
    p.id,
    v.score
FROM (
    VALUES
        ('Исландия', 'Róa', 7),
        ('Польша', 'Gaja', 10),
        ('Словения', 'How Much Time Do We Have Left', 6),
        ('Эстония', 'Espresso macchiato', 10),
        ('Испания', 'Esa Diva', 6),
        ('Украина', 'Bird of Pray', 6),
        ('Швеция', 'Bara Bada Bastu', 10),
        ('Португалия', 'Deslocado', 6),
        ('Норвегия', 'Lighter', 8),
        ('Бельгия', 'Strobe Lights', 7),
        ('Италия', 'Volevo Essere Un Duro', 7),
        ('Азербайджан', 'Run With U', 9),
        ('Сан-Марино', 'Tutta l''Italia', 7),
        ('Албания', 'Zjerm', 9),
        ('Нидерланды', 'C''est La Vie', 8),
        ('Хорватия', 'Poison Cake', 8),
        ('Швейцария', 'Voyage', 6),
        ('Кипр', 'Shh', 8)
) AS v(country_name, song, score)
JOIN countries c
    ON c.name_ru = v.country_name
JOIN performance p
    ON p.country_id = c.id
   AND p.song = v.song
JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'first-semifinal'
   AND ct.year = 2025;

INSERT INTO scores (user_id, performance_id, score)
SELECT
    v.user_id::uuid,
    p.id,
    v.score
FROM (
    VALUES
        -- Австралия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Австралия', 'Milkshake Man', 10),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Австралия', 'Milkshake Man', 7),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Австралия', 'Milkshake Man', 6),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Австралия', 'Milkshake Man', 8),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Австралия', 'Milkshake Man', 10),

        -- Черногория
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Черногория', 'Dobrodošli', 3),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Черногория', 'Dobrodošli', 9),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Черногория', 'Dobrodošli', 6),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Черногория', 'Dobrodošli', 5),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Черногория', 'Dobrodošli', 5),

        -- Ирландия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Ирландия', 'Laika Party', 4),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Ирландия', 'Laika Party', 8),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Ирландия', 'Laika Party', 4),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Ирландия', 'Laika Party', 4),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Ирландия', 'Laika Party', 3),

        -- Латвия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Латвия', 'Bur Man Laimi', 3),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Латвия', 'Bur Man Laimi', 6),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Латвия', 'Bur Man Laimi', 8),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Латвия', 'Bur Man Laimi', 8),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Латвия', 'Bur Man Laimi', 9),

        -- Армения
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Армения', 'SURVIVOR', 10),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Армения', 'SURVIVOR', 10),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Армения', 'SURVIVOR', 10),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Армения', 'SURVIVOR', 10),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Армения', 'SURVIVOR', 10),

        -- Австрия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Австрия', 'Wasted Love', 9),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Австрия', 'Wasted Love', 10),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Австрия', 'Wasted Love', 7.5),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Австрия', 'Wasted Love', 7.5),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Австрия', 'Wasted Love', 10),

        -- Великобритания
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Великобритания', 'What The Hell Just Happened?', 1),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Великобритания', 'What The Hell Just Happened?', 6),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Великобритания', 'What The Hell Just Happened?', 6),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Великобритания', 'What The Hell Just Happened?', 6.5),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Великобритания', 'What The Hell Just Happened?', 3),

        -- Греция
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Греция', 'Asteromáta', 5),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Греция', 'Asteromáta', 6),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Греция', 'Asteromáta', 5),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Греция', 'Asteromáta', 5),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Греция', 'Asteromáta', 5),

        -- Литва
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Литва', 'Tavo Akys', 4),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Литва', 'Tavo Akys', 9),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Литва', 'Tavo Akys', 5),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Литва', 'Tavo Akys', 5),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Литва', 'Tavo Akys', 8),

        -- Мальта
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Мальта', 'SERVING', 7),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Мальта', 'SERVING', 8),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Мальта', 'SERVING', 8),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Мальта', 'SERVING', 8),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Мальта', 'SERVING', 9),

        -- Грузия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Грузия', 'Свобода', 6),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Грузия', 'Свобода', 7),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Грузия', 'Свобода', 5),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Грузия', 'Свобода', 6),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Грузия', 'Свобода', 6.5),

        -- Франция
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Франция', 'Маман', 10),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Франция', 'Маман', 9),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Франция', 'Маман', 10),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Франция', 'Маман', 10),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Франция', 'Маман', 9),

        -- Дания
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Дания', 'Галлюцинация', 5),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Дания', 'Галлюцинация', 7),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Дания', 'Галлюцинация', 4),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Дания', 'Галлюцинация', 4),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Дания', 'Галлюцинация', 5),

        -- Чехия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Чехия', 'Kiss Kiss Goodbye', 8),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Чехия', 'Kiss Kiss Goodbye', 9),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Чехия', 'Kiss Kiss Goodbye', 9),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Чехия', 'Kiss Kiss Goodbye', 9),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Чехия', 'Kiss Kiss Goodbye', 9),

        -- Люксембург
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Люксембург', 'La Poupée Monte Le Son', 5),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Люксембург', 'La Poupée Monte Le Son', 6),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Люксембург', 'La Poupée Monte Le Son', 6),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Люксембург', 'La Poupée Monte Le Son', 5),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Люксембург', 'La Poupée Monte Le Son', 6),

        -- Израиль
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Израиль', 'New Day Will Rise', 5),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Израиль', 'New Day Will Rise', 9),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Израиль', 'New Day Will Rise', 6),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Израиль', 'New Day Will Rise', 6),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Израиль', 'New Day Will Rise', 5),

        -- Германия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Германия', 'Baller', 3),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Германия', 'Baller', 5),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Германия', 'Baller', 10),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Германия', 'Baller', 10),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Германия', 'Baller', 5),

        -- Сербия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Сербия', 'Mila', 1),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Сербия', 'Mila', 5),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Сербия', 'Mila', 3),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Сербия', 'Mila', 3),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Сербия', 'Mila', 4),

        -- Финляндия
        ('85cfc6c3-5549-45cb-990f-86b29e0826c9', 'Финляндия', 'ICH KOMME', 2),
        ('7c6bc8f1-1d53-4840-af70-f904a52745cc', 'Финляндия', 'ICH KOMME', 10),
        ('3324ff40-277f-486c-bbca-dbb9bc2e47ca', 'Финляндия', 'ICH KOMME', 4),
        ('48cfb5d3-4efb-403b-8599-67255ecf125c', 'Финляндия', 'ICH KOMME', 4),
        ('e0b9b315-3c6b-42ef-b836-0e0b2a1db53b', 'Финляндия', 'ICH KOMME', 7)

) AS v(user_id, country_name, song, score)
JOIN countries c
    ON c.name_ru = v.country_name
JOIN performance p
    ON p.country_id = c.id
   AND p.song = v.song
JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'second-semifinal'
   AND ct.year = 2025;


INSERT INTO scores (user_id, performance_id, score, comment)
SELECT
    '85cfc6c3-5549-45cb-990f-86b29e0826c9',
    p.id,
    v.score,
    v.comment
FROM (
    VALUES
        ('Норвегия', 10, 'но вей но вей!'),
        ('Люксембург', 6, 'веселенькая'),
        ('Эстония', 10, 'ми аморе'),
        ('Израиль', 7, 'припев прик куплеты хуйлач'),
        ('Литва', 2, 'ванечка напился и орет'),
        ('Испания', 8, 'до слоумо далеко'),
        ('Украина', 1, '(сугубо за песню)'),
        ('Великобритания', 1, 'zero points'),
        ('Австрия', 10, ''),
        ('Исландия', 6, 'прикольно моментами'),
        ('Латвия', 6, 'впрынцыпе'),
        ('Нидерланды', 7, 'весели негрик'),
        ('Финляндия', 8, 'переобулся (микрофон убил)'),
        ('Италия', 9, 'язык вытянул'),
        ('Польша', 6, ''),
        ('Германия', 4, 'за ла ла ла только угроблю😑'),
        ('Греция', 3, 'мимо'),
        ('Армения', 10, 'голый торсик зарешал'),
        ('Швейцария', 9, 'за последние 30 сек песни'),
        ('Мальта', 9, 'Ж - харизматичная'),
        ('Португалия', 1, 'хуйлач'),
        ('Дания', 4, 'дефолт просто'),
        ('Швеция', 10, 'уооо уоо уоооо'),
        ('Франция', 10, 'топ1'),
        ('Сан-Марино', 2, 'хуйлач'),
        ('Албания', 9, '')
) AS v(country_name, score, comment)

JOIN countries c
    ON c.name_ru = v.country_name

JOIN performance p
    ON p.country_id = c.id

JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'final'
   AND ct.year = 2025;

INSERT INTO scores (user_id, performance_id, score, comment)
SELECT
    '7c6bc8f1-1d53-4840-af70-f904a52745cc',
    p.id,
    v.score,
    v.comment
FROM (
    VALUES
        ('Норвегия', 8, '1 пидр'),
        ('Люксембург', 5, 'хзнеочем'),
        ('Эстония', 10, 'просто десять😂'),
        ('Израиль', 9, ''),
        ('Литва', 10, 'имба подвальная'),
        ('Испания', 3, 'просто говно😂'),
        ('Украина', 1, '(сугубо за песню)'),
        ('Великобритания', 2, ''),
        ('Австрия', 10, ''),
        ('Исландия', 7, '2 пидора'),
        ('Латвия', 6, 'хазе говно'),
        ('Нидерланды', 9, 'брок2'),
        ('Финляндия', 9, 'пэм2?'),
        ('Италия', 7, 'фу нах'),
        ('Польша', 10, 'префаером'),
        ('Германия', 2, ''),
        ('Греция', 7, ''),
        ('Армения', 10, ''),
        ('Швейцария', 2, 'скушно'),
        ('Мальта', 8, ''),
        ('Португалия', 6, 'переобулась'),
        ('Дания', 6, 'база'),
        ('Швеция', 10, ''),
        ('Франция', 10, ''),
        ('Сан-Марино', 6, 'веселый хуйлач'),
        ('Албания', 10, 'мужик краш')
) AS v(country_name, score, comment)

JOIN countries c
    ON c.name_ru = v.country_name

JOIN performance p
    ON p.country_id = c.id

JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'final'
   AND ct.year = 2025;

INSERT INTO scores (user_id, performance_id, score, comment)
SELECT
    '2ca8c8cc-7daf-42ba-a24d-b97efb8a6940',
    p.id,
    v.score,
    v.comment
FROM (
    VALUES
        ('Норвегия', 8.5, 'гут вери гут'),
        ('Люксембург', 4.5, 'хуита'),
        ('Эстония', 7, 'эспрэсо'),
        ('Израиль', 5, 'маник дороже евро'),
        ('Литва', 3, 'я срать'),
        ('Испания', 2.5, 'похожа на бейонсе, на всякий 2 поставлю, прости меня бейонсе'),
        ('Украина', 2, 'палянiця'),
        ('Великобритания', 3, 'ноу коментс'),
        ('Австрия', 10, 'WAAAAAASTED LOOOOVE'),
        ('Исландия', 5, 'Космонавты'),
        ('Латвия', 4.5, 'вызвали духа'),
        ('Нидерланды', 6.5, 'лалалала'),
        ('Финляндия', 3, 'за бубсы'),
        ('Италия', 5.5, 'Пьеро'),
        ('Польша', 4, 'ниче не пон'),
        ('Германия', 8, 'МАЙНЕ ЛИБЕ'),
        ('Греция', 4, 'пародия на Джамалу'),
        ('Армения', 9, 'и тебе тхэнкью'),
        ('Швейцария', 5.5, 'воящ'),
        ('Мальта', 7, 'бьюти бленде'),
        ('Португалия', 4, 'фортепиянист'),
        ('Дания', 3.5, 'не'),
        ('Швеция', 9.5, 'ща бы в баню и такую сосиску'),
        ('Франция', 9, 'моооя милааая МАМООО'),
        ('Сан-Марино', 1.5, 'просидела в туалете'),
        ('Албания', 3.5, 'скушно')
) AS v(country_name, score, comment)

JOIN countries c
    ON c.name_ru = v.country_name

JOIN performance p
    ON p.country_id = c.id

JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'final'
   AND ct.year = 2025;


INSERT INTO scores (user_id, performance_id, score, comment)
SELECT
    '48cfb5d3-4efb-403b-8599-67255ecf125c',
    p.id,
    v.score,
    v.comment
FROM (
    VALUES
        ('Норвегия', 9, 'ноу вэев из 10'),
        ('Люксембург', 5, 'барби из 10'),
        ('Эстония', 10, 'люблю кофе'),
        ('Израиль', 6, 'люстра из театра'),
        ('Литва', 4, 'депрессия('),
        ('Испания', 3, 'с шляпой было лучше'),
        ('Украина', 4, 'слабо'),
        ('Великобритания', 4, 'принцессы дисней'),
        ('Австрия', 10, 'уплываю на гаити'),
        ('Исландия', 5, 'майнкафт моя жизнь'),
        ('Латвия', 7.5, 'лагуны'),
        ('Нидерланды', 9, 'ляля класс'),
        ('Финляндия', 4, 'сиськи класс'),
        ('Италия', 10, 'песня вау'),
        ('Польша', 3, 'смыв'),
        ('Германия', 9.5, 'нанана'),
        ('Греция', 4, 'плюсую с полиной'),
        ('Армения', 10, 'горяч.'),
        ('Швейцария', 9, 'девочка с картинки'),
        ('Мальта', 8, 'сервинг'),
        ('Португалия', 6.5, ''),
        ('Дания', 3, 'хуйня'),
        ('Швеция', 10, 'банька парилка ебля и горилка'),
        ('Франция', 9.7, 'очень хорош'),
        ('Сан-Марино', 3, 'не оч'),
        ('Албания', 6, 'нормис и красный себастьян')
) AS v(country_name, score, comment)

JOIN countries c
    ON c.name_ru = v.country_name

JOIN performance p
    ON p.country_id = c.id

JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'final'
   AND ct.year = 2025;

INSERT INTO scores (user_id, performance_id, score, comment)
SELECT
    '3324ff40-277f-486c-bbca-dbb9bc2e47ca',
    p.id,
    v.score,
    v.comment
FROM (
    VALUES
        ('Норвегия', 9, 'рыцарей из 10'),
        ('Люксембург', 5, 'хочу такой пиджак'),
        ('Эстония', 10, 'это я выбежала'),
        ('Израиль', 2, '*блюющий стикер*'),
        ('Литва', 3, 'ТРАВА ТРАВА ТРАВА'),
        ('Испания', 3, 'в начале показалось что буде хорошо'),
        ('Украина', 3, 'слабо  для украины'),
        ('Великобритания', 3, 'зироооу пойнт'),
        ('Австрия', 9, 'витасик'),
        ('Исландия', 4, 'ваеб'),
        ('Латвия', 8, 'ваканда в сердечке'),
        ('Нидерланды', 9, 'негров из 10'),
        ('Финляндия', 4, 'сиськи класс'),
        ('Италия', 10, 'мой дэвид боучик'),
        ('Польша', 4, 'за 52 года'),
        ('Германия', 10, 'нананананананананнана'),
        ('Греция', 4, 'согл с полиной'),
        ('Армения', 10, 'слушабельно'),
        ('Швейцария', 8, 'мона лиза'),
        ('Мальта', 9, 'Я ПИЗДАЯ'),
        ('Португалия', 6, 'не хватает искры'),
        ('Дания', 2, 'ни о чем'),
        ('Швеция', 10, 'когда на дачу?  когда Давид или Полина дадут добро'),
        ('Франция', 8, 'да.'),
        ('Сан-Марино', 1, 'ИМБИЩЕ!'),
        ('Албания', 2, 'хунч')
) AS v(country_name, score, comment)

JOIN countries c
    ON c.name_ru = v.country_name

JOIN performance p
    ON p.country_id = c.id

JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'final'
   AND ct.year = 2025;

INSERT INTO scores (user_id, performance_id, score, comment)
SELECT
    'e0b9b315-3c6b-42ef-b836-0e0b2a1db53b',
    p.id,
    v.score,
    v.comment
FROM (
    VALUES
        ('Норвегия', 10, 'Сыночка давай'),
        ('Люксембург', 3, 'Я не игрушка оуоуоуо'),
        ('Эстония', 10, 'Кофе мой друууг'),
        ('Израиль', 2, 'Дорогая, ты не Эльза и здесь не твой дом'),
        ('Литва', 6, 'Не мешайте ему, он в эдите'),
        ('Испания', 4, 'Вива ля дива, вива Виктория, афродитаааа. ,p.s. Киркоров'),
        ('Украина', 6, 'Литл берд…. ЭЭ🦅 немного литтле пёрд'),
        ('Великобритания', 4, 'Ладно, пойдем по нашим песням: одна река была как белый деееень…'),
        ('Австрия', 10, 'Наш корабль идет ко дну… помоги мне я ща тут все обосру...'),
        ('Исландия', 5, 'Клетки моего мозга на экзамене:'),
        ('Латвия', 8, 'А вот это уже Винкс- Сиреникс!'),
        ('Нидерланды', 9, 'Се ля ви, Мон ами, ведь можно жить в блеске! (Ладно, он тоже в эдите)'),
        ('Финляндия', 4, 'Присоединяюсь к комментариям про сиськи. (•Y•)'),
        ('Италия', 7, 'Как из мультика Тима Бертона'),
        ('Польша', 4, 'Хлооопай ресницами и взлетай… сниться ток не надо пж..'),
        ('Германия', 7, 'Светка ночью на даче всех будит сходить в туалет с ней. ХАХАХАХАХАХ ДА.'),
        ('Греция', 4, 'Закос на маму Терезу'),
        ('Армения', 10, 'Бежит груша - нельзя скушать'),
        ('Швейцария', 9, 'Монеточка?'),
        ('Мальта', 9, 'Плачь Европа, ведь у меня самая красивая опа'),
        ('Португалия', 5, 'Мне даже не за что зацепиться и что-то придумать('),
        ('Дания', 3, 'В синем море, в белой пене'),
        ('Швеция', 10, 'Ребят, от души, все в баню, нас пригласили мужички'),
        ('Франция', 10, 'Не обращайте внимание, это со Светы песок сыпется'),
        ('Сан-Марино', 3, 'Ладно, тут он 100% занял место Хорвата или Кипра в финале. Не прощаем.'),
        ('Албания', 7, 'Вот с кем уехала мама из Папиных Дочек…')
) AS v(country_name, score, comment)

JOIN countries c
    ON c.name_ru = v.country_name

JOIN performance p
    ON p.country_id = c.id

JOIN contests ct
    ON ct.id = p.contest_id
   AND ct.type = 'final'
   AND ct.year = 2025;

UPDATE performance p
SET youtube_link = v.youtube_link
FROM (
    VALUES
        ('Норвегия', 'https://youtu.be/gQOGxx6Fk9k?si=dU47S6pIZ_SKEUB1'),
        ('Люксембург', 'https://youtu.be/GT7ZZBCscUg?si=b2cmQp5uVX6VbP7U'),
        ('Эстония', 'https://youtu.be/9b9Z5HSCXOI?si=aePOzXqNjJq8Eaj_'),
        ('Израиль', 'https://youtu.be/_7zHp51j2WM?si=0CzRvs31nG0Z5PS_'),
        ('Литва', 'https://youtu.be/3F6bwWGhm_s?si=d76g7RgKL1gIvDru'),
        ('Испания', 'https://youtu.be/IEKSa9FVLqA?si=4rTop1hkuoCy4Tz8'),
        ('Украина', 'https://youtu.be/-DG0l8sSNJM?si=dNzn5vQrR8J4UgFl'),
        ('Великобритания', 'https://youtu.be/Ur5qRh0BaHk?si=sUnMYPAfaoZMc-lK'),
        ('Австрия', 'https://youtu.be/onOex2WXjbA?si=GZQnQxBy9aj3yENI'),
        ('Исландия', 'https://youtu.be/c73Lx1QUZZA?si=GSFd02nm55FeTZPn'),
        ('Латвия', 'https://youtu.be/nkvcMe3NiQ0?si=QsUpYpckj_IVwIg7'),
        ('Нидерланды', 'https://youtu.be/LiTQVJwxvfE?si=SY0HTwitQj1reXtq'),
        ('Финляндия', 'https://youtu.be/V3vbVd1ynnk?si=EVU-QuyfXZ3-dbM8'),
        ('Италия', 'https://youtu.be/Vlu5XXDwHos?si=SjwHcs12X2lmsAWW'),
        ('Польша', 'https://youtu.be/eg5RtEX1zJ0?si=nWI09BEggO4twP0A'),
        ('Германия', 'https://youtu.be/3rrWZ6cldsA?si=JiS6Q8Ts3RSRKtB9'),
        ('Греция', 'https://youtu.be/1qbWRl6h6to?si=gbWmFDvVva5hvvZK'),
        ('Армения', 'https://youtu.be/qHkZWLld-pw?si=4kTTzAXXEfk3SdRB'),
        ('Швейцария', 'https://youtu.be/5TMc6HzimQo?si=tWF-Q1KgIjz0Yu3v'),
        ('Мальта', 'https://youtu.be/povnGP6k0sI?si=hUUpxlPKZXm0tXcp'),
        ('Португалия', 'https://youtu.be/waInyqBwSo0?si=em0sWWGtnML9-hg8'),
        ('Дания', 'https://youtu.be/B3BdsYDnS8M?si=Q4BITBjVeSvun7vD'),
        ('Швеция', 'https://youtu.be/WSh7U3m9KgA?si=m8XafeuF-y48LgQ3'),
        ('Франция', 'https://youtu.be/jhqJY0ll1Wo?si=A67JEzsh_StTlbyC'),
        ('Сан-Марино', 'https://youtu.be/hq6XIRKmA2A?si=KyVs3jvSaNm_6vSZ'),
        ('Албания', 'https://youtu.be/xfn6ssOf_zU?si=ygeEsqURW3FvYE8u'),
        ('Кипр', 'https://youtu.be/egPAiAuC57k?si=vrJhqybRiewXJ_VL'),
        ('Хорватия', 'https://youtu.be/jzK4D_gfRjQ?si=3mFGGOxaCuoo2TrH'),
        ('Бельгия', 'https://youtu.be/fl4LaADiLBY?si=AVpe2_EkEjf2sjnM'),
        ('Азербайджан', 'https://youtu.be/wk1CUjaRKyo?si=UkFJu_NxcvUgXir4'),
        ('Словения', 'https://youtu.be/Jbs9WlvIkg0?si=YBrWHmHpeBN5QNAe'),
        ('Чехия', 'https://youtu.be/hdxna1DC7yo?si=H6qaR0jNmvHXqGYB'),
        ('Австралия', 'https://youtu.be/EJ0RdIU_G8g?si=DoSLtDmUUW4BLnQ1'),
        ('Грузия', 'https://youtu.be/jphJoo-CNtU?si=HruQWC4Y1fVDkKhF'),
        ('Черногория', 'https://youtu.be/L9MNHACTvT0?si=fwv4HxSuOccy8CrE'),
        ('Ирландия', 'https://youtu.be/3MB628Kanzo?si=zwNh15uCDtQ66GzK'),
        ('Сербия', 'https://youtu.be/WlCoZ0UOXoY?si=p3QrlMg75GXfiKLo')
) AS v(country_name, youtube_link)
JOIN countries c
    ON c.name_ru = v.country_name
JOIN performance p2
    ON p2.country_id = c.id
JOIN contests ct
    ON ct.id = p2.contest_id
   AND ct.year = 2025
WHERE p.id = p2.id;