create table users (
	id uuid primary key,
    username text not null unique,
    created_at timestamp not null,
    password text not null
);

CREATE TYPE contest_type AS ENUM ('final', 'first-semifinal', 'second-semifinal');
create table contests (
	id uuid primary key,
	type contest_type not null,
	year int not null 
);

create table contest_participants (
	user_id uuid REFERENCES users(id),
	contest_id uuid REFERENCES contests(id)
);

create table countries (
	id uuid primary key, 
	name_ru text not null, 
	flag_emogi text not null
);

create table performance (
	id uuid primary key, 
	contest_id uuid REFERENCES contests(id),
	country_id uuid REFERENCES countries(id),
	number int not null,
	youtube_link text not null,
	artist text not null,
	song text not null
);

create table scores (
	user_id uuid REFERENCES users(id),
	performance_id uuid REFERENCES performance(id),
	score int not null,
	comment text,

	unique (user_id, performance_id)
);

INSERT INTO countries (id, name_ru, flag_emogi) VALUES

(gen_random_uuid(), 'Молдова', '🇲🇩'),
(gen_random_uuid(), 'Швеция', '🇸🇪'),
(gen_random_uuid(), 'Хорватия', '🇭🇷'),
(gen_random_uuid(), 'Греция', '🇬🇷'),
(gen_random_uuid(), 'Португалия', '🇵🇹'),
(gen_random_uuid(), 'Грузия', '🇬🇪'),
(gen_random_uuid(), 'Италия', '🇮🇹'),
(gen_random_uuid(), 'Финляндия', '🇫🇮'),
(gen_random_uuid(), 'Черногория', '🇲🇪'),
(gen_random_uuid(), 'Эстония', '🇪🇪'),
(gen_random_uuid(), 'Израиль', '🇮🇱'),
(gen_random_uuid(), 'Германия', '🇩🇪'),
(gen_random_uuid(), 'Бельгия', '🇧🇪'),
(gen_random_uuid(), 'Литва', '🇱🇹'),
(gen_random_uuid(), 'Сан-Марино', '🇸🇲'),
(gen_random_uuid(), 'Польша', '🇵🇱'),
(gen_random_uuid(), 'Сербия', '🇷🇸'),
(gen_random_uuid(), 'Болгария', '🇧🇬'),
(gen_random_uuid(), 'Азербайджан', '🇦🇿'),
(gen_random_uuid(), 'Румыния', '🇷🇴'),
(gen_random_uuid(), 'Люксембург', '🇱🇺'),
(gen_random_uuid(), 'Чехия', '🇨🇿'),
(gen_random_uuid(), 'Франция', '🇫🇷'),
(gen_random_uuid(), 'Армения', '🇦🇲'),
(gen_random_uuid(), 'Швейцария', '🇨🇭'),
(gen_random_uuid(), 'Кипр', '🇨🇾'),
(gen_random_uuid(), 'Австрия', '🇦🇹'),
(gen_random_uuid(), 'Латвия', '🇱🇻'),
(gen_random_uuid(), 'Дания', '🇩🇰'),
(gen_random_uuid(), 'Австралия', '🇦🇺'),
(gen_random_uuid(), 'Украина', '🇺🇦'),
(gen_random_uuid(), 'Великобритания', '🇬🇧'),
(gen_random_uuid(), 'Албания', '🇦🇱'),
(gen_random_uuid(), 'Мальта', '🇲🇹'),
(gen_random_uuid(), 'Норвегия', '🇳🇴');

INSERT INTO contests (id, type, year)
VALUES
(gen_random_uuid(), 'first-semifinal', 2026),
(gen_random_uuid(), 'second-semifinal', 2026);

WITH c AS (
    SELECT id, type FROM contests WHERE year = 2026
)

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

-- FIRST SEMIFINAL
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Молдова'), 1,  'https://youtube.com/x1', 'Satoshi', 'Viva, Moldova!'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Швеция'), 2, 'https://youtube.com/x2', 'FELICIA', 'My System'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Хорватия'), 3, 'https://youtube.com/x3', 'LELEK', 'Andromeda'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Греция'), 4, 'https://youtube.com/x4', 'Akylas', 'Ferto'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Португалия'), 5, 'https://youtube.com/x5', 'Bandidos do Cante', 'Rosa'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Грузия'), 6, 'https://youtube.com/x6', 'Bzikebi', 'On Replay'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Италия'), 7, 'https://youtube.com/x7', 'Sal Da Vinci', 'Per sempre sì'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Финляндия'), 8, 'https://youtube.com/x8', 'Linda Lampenius x Pete Parkkonen', 'Liekinheitin'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Черногория'), 9, 'https://youtube.com/x9', 'Tamara Živković', 'Nova Zora'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Эстония'), 10, 'https://youtube.com/x10', 'Vanilla Ninja', 'Too Epic To Be True'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Израиль'), 11, 'https://youtube.com/x11', 'Noam Bettan', 'Michelle'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Германия'), 12, 'https://youtube.com/x12', 'Sarah Engels', 'Fire'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Бельгия'), 13, 'https://youtube.com/x13', 'ESSYLA', 'Dancing on the Ice'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Литва'), 14, 'https://youtube.com/x14', 'Lion Ceccah', 'Sólo Quiero Más'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Сан-Марино'), 15, 'https://youtube.com/x15', 'Senhit feat. Boy George', 'Superstar'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Польша'), 16, 'https://youtube.com/x16', 'ALICJA', 'Pray'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='first-semifinal'), (SELECT id FROM countries WHERE name_ru='Сербия'), 17, 'https://youtube.com/x17', 'LAVINA', 'Kraj mene'),

-- SECOND SEMIFINAL
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Болгария'), 1, 'https://youtube.com/x18', 'DARA', 'Bangaranga'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Азербайджан'), 2, 'https://youtube.com/x19', 'JIVA', 'Just Go'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Румыния'), 3, 'https://youtube.com/x20', 'Alexandra Căpitănescu', 'Choke Me'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Люксембург'), 4, 'https://youtube.com/x21', 'Eva Marija', 'Mother Nature'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Чехия'), 5, 'https://youtube.com/x22', 'Daniel Zizka', 'CROSSROADS'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Франция'), 6, 'https://youtube.com/x23', 'Monroe', 'Regarde !'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Армения'), 7, 'https://youtube.com/x24', 'SIMÓN', 'Paloma Rumba'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Швейцария'), 8, 'https://youtube.com/x25', 'Veronica Fusaro', 'Alice'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Кипр'), 9, 'https://youtube.com/x26', 'Antigoni', 'JALLA'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Австрия'), 10, 'https://youtube.com/x27', 'COSMÓ', 'Tanzschein'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Латвия'), 11, 'https://youtube.com/x28', 'Atvara', 'Ēnā'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Дания'), 12, 'https://youtube.com/x29', 'Søren Torpegaard Lund', 'Før Vi Går Hjem'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Австралия'), 13, 'https://youtube.com/x30', 'Delta Goodrem', 'Eclipse'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Украина'), 14, 'https://youtube.com/x31', 'LELÉKA', 'Ridnym'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Великобритания'), 15, 'https://youtube.com/x32', 'LOOK MUM NO COMPUTER', 'Eins, Zwei, Drei'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Албания'), 16, 'https://youtube.com/x33', 'Alis', 'Nân'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Мальта'), 17, 'https://youtube.com/x34', 'AIDAN', 'Bella'),
(gen_random_uuid(), (SELECT id FROM c WHERE type='second-semifinal'), (SELECT id FROM countries WHERE name_ru='Норвегия'), 18, 'https://youtube.com/x35', 'JONAS LOVV', 'YA YA YA');