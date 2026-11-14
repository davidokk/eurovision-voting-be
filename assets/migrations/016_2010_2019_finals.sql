-- Eurovision Song Contest 2010-2019 — Grand Finals
-- Running order, artists, songs and official ESC YouTube live performances.

-- Extra countries for 2010-2019 finals
INSERT INTO countries (id, name_ru, flag_emogi)
SELECT gen_random_uuid(), 'Босния и Герцеговина', '🇧🇦'
WHERE NOT EXISTS (SELECT 1 FROM countries WHERE name_ru = 'Босния и Герцеговина');

INSERT INTO countries (id, name_ru, flag_emogi)
SELECT gen_random_uuid(), 'Беларусь', '🇧🇾'
WHERE NOT EXISTS (SELECT 1 FROM countries WHERE name_ru = 'Беларусь');

INSERT INTO countries (id, name_ru, flag_emogi)
SELECT gen_random_uuid(), 'Венгрия', '🇭🇺'
WHERE NOT EXISTS (SELECT 1 FROM countries WHERE name_ru = 'Венгрия');

INSERT INTO countries (id, name_ru, flag_emogi)
SELECT gen_random_uuid(), 'Македония', '🇲🇰'
WHERE NOT EXISTS (SELECT 1 FROM countries WHERE name_ru = 'Македония');

INSERT INTO countries (id, name_ru, flag_emogi)
SELECT gen_random_uuid(), 'Северная Македония', '🇲🇰'
WHERE NOT EXISTS (SELECT 1 FROM countries WHERE name_ru = 'Северная Македония');

INSERT INTO countries (id, name_ru, flag_emogi)
SELECT gen_random_uuid(), 'Турция', '🇹🇷'
WHERE NOT EXISTS (SELECT 1 FROM countries WHERE name_ru = 'Турция');

-- 2010 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2010,
    '2010-05-29 21:00:00+02',
    '2010-05-30 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 1, 'https://www.youtube.com/watch?v=SOKGZb38e5c', 'Safura', 'Drip Drop'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Испания'), 2, 'https://www.youtube.com/watch?v=6gltZu8lsFc', 'Daniel Diges', 'Algo pequeñito'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 3, 'https://www.youtube.com/watch?v=9w527tADNpc', 'Didrik Solli-Tangen', 'My Heart Is Yours'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Молдова'), 4, 'https://www.youtube.com/watch?v=pHXDMe6QV-U', 'SunStroke Project and Olia Tira', 'Run Away'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 5, 'https://www.youtube.com/watch?v=ZGNtoR2hom4', 'Jon Lilygreen and the Islanders', 'Life Looks Better in Spring'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Босния и Герцеговина'), 6, 'https://www.youtube.com/watch?v=IPBo9vaN9Vc', 'Vukašin Brajić', 'Thunder and Lightning'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Бельгия'), 7, 'https://www.youtube.com/watch?v=lTsnwGrPKEI', 'Tom Dice', 'Me and My Guitar'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 8, 'https://www.youtube.com/watch?v=4pVEII2pmso', 'Milan Stanković', 'Ovo je Balkan'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Беларусь'), 9, 'https://www.youtube.com/watch?v=MvK2KU2zfzA', '3+2 feat. Robert Wells', 'Butterflies'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Ирландия'), 10, 'https://www.youtube.com/watch?v=6SHg8gS-02Q', 'Niamh Kavanagh', 'It''s for You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Греция'), 11, 'https://www.youtube.com/watch?v=wfqOrMlwOjs', 'Giorgos Alkaios and Friends', 'Opa'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 12, 'https://www.youtube.com/watch?v=Nro3hlemmho', 'Josh Dubovie', 'That Sounds Good to Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Грузия'), 13, 'https://www.youtube.com/watch?v=y5yTMGxFiPw', 'Sofia Nizharadze', 'Shine'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Турция'), 14, 'https://www.youtube.com/watch?v=HB_GnnhNz-8', 'Manga', 'We Could Be the Same'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Албания'), 15, 'https://www.youtube.com/watch?v=idE_32NUJbo', 'Juliana Pasha', 'It''s All About You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Исландия'), 16, 'https://www.youtube.com/watch?v=qXLqw4UpJ1w', 'Hera Björk', 'Je ne sais quoi'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Украина'), 17, 'https://www.youtube.com/watch?v=20CPPV_uwmE', 'Alyosha', 'Sweet People'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Франция'), 18, 'https://www.youtube.com/watch?v=-NlvQ6Ln9fc', 'Jessy Matador', 'Allez Ola Olé'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Румыния'), 19, 'https://www.youtube.com/watch?v=J9EtMZXeQZw', 'Paula Seling and Ovi', 'Playing with Fire'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Россия'), 20, 'https://www.youtube.com/watch?v=nFEwKjpDD2c', 'Peter Nalitch and Friends', 'Lost and Forgotten'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Армения'), 21, 'https://www.youtube.com/watch?v=A60atIirQAw', 'Eva Rivas', 'Apricot Stone'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Германия'), 22, 'https://www.youtube.com/watch?v=7pL9vdpSvnY', 'Lena', 'Satellite'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Португалия'), 23, 'https://www.youtube.com/watch?v=TvmnvtL_DcQ', 'Filipa Azevedo', 'Há dias assim'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Израиль'), 24, 'https://www.youtube.com/watch?v=Tj0qaNnLk9M', 'Harel Skaat', 'Milim'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2010), (SELECT id FROM countries WHERE name_ru = 'Дания'), 25, 'https://www.youtube.com/watch?v=Px2UsvO1lCg', 'Chanée and N''evergreen', 'In a Moment like This');

-- 2011 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2011,
    '2011-05-14 21:00:00+02',
    '2011-05-15 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Финляндия'), 1, 'https://www.youtube.com/watch?v=rQ98uBQxbBo', 'Paradise Oskar', 'Da Da Dam'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Босния и Герцеговина'), 2, 'https://www.youtube.com/watch?v=mBg3coarF_8', 'Dino Merlin', 'Love in Rewind'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Дания'), 3, 'https://www.youtube.com/watch?v=9DcYb8kt4mc', 'A Friend in London', 'New Tomorrow'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Литва'), 4, 'https://www.youtube.com/watch?v=D4Tvc9S_1BU', 'Evelina Sašenko', 'C''est ma vie'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 5, 'https://www.youtube.com/watch?v=kwH7Pw25kNk', 'Kati Wolf', 'What About My Dreams?'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Ирландия'), 6, 'https://www.youtube.com/watch?v=rEyd-DVLq9M', 'Jedward', 'Lipstick'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 7, 'https://www.youtube.com/watch?v=Zhn4DFHhM3U', 'Eric Saade', 'Popular'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Эстония'), 8, 'https://www.youtube.com/watch?v=vOIXQskfy6o', 'Getter Jaani', 'Rockefeller Street'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Греция'), 9, 'https://www.youtube.com/watch?v=FdZzklXoAOg', 'Loukas Yorkas feat. Stereo Mike', 'Watch My Dance'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Россия'), 10, 'https://www.youtube.com/watch?v=C-3ohaRN3MM', 'Alexey Vorobyov', 'Get You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Франция'), 11, 'https://www.youtube.com/watch?v=Sri0LomK0oM', 'Amaury Vassili', 'Sognu'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Италия'), 12, 'https://www.youtube.com/watch?v=TE0uNLp3LuU', 'Raphael Gualazzi', 'Madness of Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Швейцария'), 13, 'https://www.youtube.com/watch?v=M03Nx6gHd90', 'Anna Rossinelli', 'In Love for a While'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 14, 'https://www.youtube.com/watch?v=y1dXMWWq6y8', 'Blue', 'I Can'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Молдова'), 15, 'https://www.youtube.com/watch?v=lyQ72wb3yr4', 'Zdob și Zdub', 'So Lucky'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Германия'), 16, 'https://www.youtube.com/watch?v=KqZprvpWlcE', 'Lena', 'Taken by a Stranger'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Румыния'), 17, 'https://www.youtube.com/watch?v=2a1Q_hpGLbs', 'Hotel FM', 'Change'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Австрия'), 18, 'https://www.youtube.com/watch?v=AySl4viIzlo', 'Nadine Beiler', 'The Secret Is Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 19, 'https://www.youtube.com/watch?v=_0tlQUW5X0U', 'Ell and Nikki', 'Running Scared'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Словения'), 20, 'https://www.youtube.com/watch?v=5vqjb2OYnYs', 'Maja Keuc', 'No One'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Исландия'), 21, 'https://www.youtube.com/watch?v=apc_qJf3nws', 'Sjonni''s Friends', 'Coming Home'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Испания'), 22, 'https://www.youtube.com/watch?v=GSOEmM0h-aE', 'Lucía Pérez', 'Que me quiten lo bailao'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Украина'), 23, 'https://www.youtube.com/watch?v=mx1wAj_oVdE', 'Mika Newton', 'Angel'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 24, 'https://www.youtube.com/watch?v=oHgXE90OSzM', 'Nina', 'Čaroban'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2011), (SELECT id FROM countries WHERE name_ru = 'Грузия'), 25, 'https://www.youtube.com/watch?v=s2EwiTBiEDU', 'Eldrine', 'One More Day');

-- 2012 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2012,
    '2012-05-26 21:00:00+02',
    '2012-05-27 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 1, 'https://www.youtube.com/watch?v=rXw4Q5jbNqQ', 'Engelbert Humperdinck', 'Love Will Set You Free'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 2, 'https://www.youtube.com/watch?v=sw85G2AlU2c', 'Compact Disco', 'Sound of Our Hearts'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Албания'), 3, 'https://www.youtube.com/watch?v=QeBL2UHhyEc', 'Rona Nishliu', 'Suus'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Литва'), 4, 'https://www.youtube.com/watch?v=x2vNoZnJxgY', 'Donny Montell', 'Love Is Blind'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Босния и Герцеговина'), 5, 'https://www.youtube.com/watch?v=81hIbZNoFU8', 'Maya Sar', 'Korake ti znam'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Россия'), 6, 'https://www.youtube.com/watch?v=BgUstrmJzyc', 'Buranovskiye Babushki', 'Party for Everybody'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Исландия'), 7, 'https://www.youtube.com/watch?v=zRqmHRbJvTo', 'Greta Salóme and Jónsi', 'Never Forget'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 8, 'https://www.youtube.com/watch?v=ex133UhxB64', 'Ivi Adamou', 'La La Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Франция'), 9, 'https://www.youtube.com/watch?v=0DZ7DbdeOeE', 'Anggun', 'Echo (You and I)'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Италия'), 10, 'https://www.youtube.com/watch?v=v0kGpDEvtbQ', 'Nina Zilli', 'L''amore è femmina (Out of Love)'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Эстония'), 11, 'https://www.youtube.com/watch?v=T7BaTBe0UD8', 'Ott Lepland', 'Kuula'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 12, 'https://www.youtube.com/watch?v=-ZaxIY1VAbM', 'Tooji', 'Stay'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 13, 'https://www.youtube.com/watch?v=yzT7O3Fnwpk', 'Sabina Babayeva', 'When the Music Dies'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Румыния'), 14, 'https://www.youtube.com/watch?v=K3ruy639kTQ', 'Mandinga', 'Zaleilah'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Дания'), 15, 'https://www.youtube.com/watch?v=U5X2r_t-KBk', 'Soluna Samay', 'Should''ve Known Better'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Греция'), 16, 'https://www.youtube.com/watch?v=NreIi8ma_PY', 'Eleftheria Eleftheriou', 'Aphrodisiac'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 17, 'https://www.youtube.com/watch?v=Pfo-8z86x80', 'Loreen', 'Euphoria'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Турция'), 18, 'https://www.youtube.com/watch?v=3Qa7_y21oOY', 'Can Bonomo', 'Love Me Back'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Испания'), 19, 'https://www.youtube.com/watch?v=U8J1b62wOao', 'Pastora Soler', 'Quédate conmigo'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Германия'), 20, 'https://www.youtube.com/watch?v=cScJPH20P3A', 'Roman Lob', 'Standing Still'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Мальта'), 21, 'https://www.youtube.com/watch?v=lUnEZp-qHxs', 'Kurt Calleja', 'This Is the Night'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Македония'), 22, 'https://www.youtube.com/watch?v=MSuCdNEy9w0', 'Kaliopi', 'Crno i belo'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Ирландия'), 23, 'https://www.youtube.com/watch?v=e1cuimKLNpU', 'Jedward', 'Waterline'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 24, 'https://www.youtube.com/watch?v=R9x9VbJzaDQ', 'Željko Joksimović', 'Nije ljubav stvar'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Украина'), 25, 'https://www.youtube.com/watch?v=sLsTn_li5d8', 'Gaitana', 'Be My Guest'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2012), (SELECT id FROM countries WHERE name_ru = 'Молдова'), 26, 'https://www.youtube.com/watch?v=vIy0EtEeLEw', 'Pasha Parfeny', 'Lăutar');

-- 2013 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2013,
    '2013-05-18 21:00:00+02',
    '2013-05-19 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Франция'), 1, 'https://www.youtube.com/watch?v=hQwiQj5Dyj8', 'Amandine Bourgeois', 'L''Enfer et moi'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Литва'), 2, 'https://www.youtube.com/watch?v=0VLhnzk_dAs', 'Andrius Pojavis', 'Something'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Молдова'), 3, 'https://www.youtube.com/watch?v=xckgLUv73Jc', 'Aliona Moon', 'O mie'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Финляндия'), 4, 'https://www.youtube.com/watch?v=dlBXOveVh7c', 'Krista Siegfrids', 'Marry Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Испания'), 5, 'https://www.youtube.com/watch?v=5o_DnM_ANF8', 'ESDM', 'Contigo hasta el final'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Бельгия'), 6, 'https://www.youtube.com/watch?v=dCl7OUFNSCA', 'Roberto Bellarosa', 'Love Kills'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Эстония'), 7, 'https://www.youtube.com/watch?v=Ooo6GvZNLhI', 'Birgit', 'Et uus saaks alguse'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Беларусь'), 8, 'https://www.youtube.com/watch?v=XKewBfUE2Mg', 'Alyona Lanskaya', 'Solayoh'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Мальта'), 9, 'https://www.youtube.com/watch?v=VK4TvZeL6c0', 'Gianluca', 'Tomorrow'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Россия'), 10, 'https://www.youtube.com/watch?v=wHNwk3Oez8U', 'Dina Garipova', 'What If'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Германия'), 11, 'https://www.youtube.com/watch?v=3HYTd7WDCPg', 'Cascada', 'Glorious'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Армения'), 12, 'https://www.youtube.com/watch?v=UGOSZ7Uufno', 'Dorians', 'Lonely Planet'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Нидерланды'), 13, 'https://www.youtube.com/watch?v=n5iazXvMw5o', 'Anouk', 'Birds'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Румыния'), 14, 'https://www.youtube.com/watch?v=VgHWFiavqjA', 'Cezar', 'It''s My Life'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 15, 'https://www.youtube.com/watch?v=-agn-ERgTXA', 'Bonnie Tyler', 'Believe in Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 16, 'https://www.youtube.com/watch?v=vtjdTPnCcu0', 'Robin Stjernberg', 'You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 17, 'https://www.youtube.com/watch?v=rOimfHq76xk', 'ByeAlex', 'Kedvesem'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Дания'), 18, 'https://www.youtube.com/watch?v=p3f9v8ebuD4', 'Emmelie de Forest', 'Only Teardrops'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Исландия'), 19, 'https://www.youtube.com/watch?v=GzE88IYzXLI', 'Eythor Ingi', 'Ég á líf'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 20, 'https://www.youtube.com/watch?v=iN3d_V7KVLE', 'Farid Mammadov', 'Hold Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Греция'), 21, 'https://www.youtube.com/watch?v=_AfrOZersNk', 'Koza Mostra feat. Agathon Iakovidis', 'Alcohol Is Free'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Украина'), 22, 'https://www.youtube.com/watch?v=BMxpv8naRd8', 'Zlata Ognevich', 'Gravity'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Италия'), 23, 'https://www.youtube.com/watch?v=S8oaxDV1q6o', 'Marco Mengoni', 'L''essenziale'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 24, 'https://www.youtube.com/watch?v=tc6a4EV63uM', 'Margaret Berger', 'I Feed You My Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Грузия'), 25, 'https://www.youtube.com/watch?v=-n3mwx0vQDE', 'Nodi Tatishvili and Sophie Gelovani', 'Waterfall'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2013), (SELECT id FROM countries WHERE name_ru = 'Ирландия'), 26, 'https://www.youtube.com/watch?v=4tgtyl_5x5Y', 'Ryan Dolan', 'Only Love Survives');

-- 2014 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2014,
    '2014-05-10 21:00:00+02',
    '2014-05-11 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Украина'), 1, 'https://www.youtube.com/watch?v=slHboKF9PIQ', 'Mariya Yaremchuk', 'Tick-Tock'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Беларусь'), 2, 'https://www.youtube.com/watch?v=tImGMW-4eHY', 'Teo', 'Cheesecake'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 3, 'https://www.youtube.com/watch?v=5mPUMPcFQAY', 'Dilara Kazimova', 'Start a Fire'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Исландия'), 4, 'https://www.youtube.com/watch?v=U2qAvGsiCmc', 'Pollapönk', 'No Prejudice'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 5, 'https://www.youtube.com/watch?v=BTMli8U17MM', 'Carl Espen', 'Silent Storm'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Румыния'), 6, 'https://www.youtube.com/watch?v=qN2h6__yJGQ', 'Paula Seling and Ovi', 'Miracle'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Армения'), 7, 'https://www.youtube.com/watch?v=oj0oOV-2fRQ', 'Aram Mp3', 'Not Alone'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Черногория'), 8, 'https://www.youtube.com/watch?v=6FJtzRSKHKg', 'Sergej Ćetković', 'Moj svijet'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Польша'), 9, 'https://www.youtube.com/watch?v=VJ920cN2HmA', 'Donatan and Cleo', 'My Słowianie – We Are Slavic'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Греция'), 10, 'https://www.youtube.com/watch?v=c9y6A1uXLn8', 'Freaky Fortune feat. RiskyKidd', 'Rise Up'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Австрия'), 11, 'https://www.youtube.com/watch?v=QRUIava4WRM', 'Conchita Wurst', 'Rise Like a Phoenix'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Германия'), 12, 'https://www.youtube.com/watch?v=BFuGH2Kp9dI', 'Elaiza', 'Is It Right'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 13, 'https://www.youtube.com/watch?v=5PQJI-3LW-8', 'Sanna Nielsen', 'Undo'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Франция'), 14, 'https://www.youtube.com/watch?v=vjkqciwP034', 'Twin Twin', 'Moustache'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Россия'), 15, 'https://www.youtube.com/watch?v=BCNN1JAgWqk', 'Tolmachevy Sisters', 'Shine'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Италия'), 16, 'https://www.youtube.com/watch?v=8W5ZlalAMV8', 'Emma', 'La mia città'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Словения'), 17, 'https://www.youtube.com/watch?v=CiwM1sY3u0s', 'Tinkara Kovač', 'Round and Round'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Финляндия'), 18, 'https://www.youtube.com/watch?v=3FDWaFfo1CU', 'Softengine', 'Something Better'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Испания'), 19, 'https://www.youtube.com/watch?v=LY4n5FUlfDQ', 'Ruth Lorenzo', 'Dancing in the Rain'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Швейцария'), 20, 'https://www.youtube.com/watch?v=DqjoM8ZlyMc', 'Sebalter', 'Hunter of Stars'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 21, 'https://www.youtube.com/watch?v=T_RxJZ3N9kE', 'András Kállay-Saunders', 'Running'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Мальта'), 22, 'https://www.youtube.com/watch?v=E_hT08vzHt8', 'Firelight', 'Coming Home'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Дания'), 23, 'https://www.youtube.com/watch?v=jnCTsJkieY4', 'Basim', 'Cliche Love Song'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Нидерланды'), 24, 'https://www.youtube.com/watch?v=4ggBPAm5XLA', 'The Common Linnets', 'Calm After the Storm'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Сан-Марино'), 25, 'https://www.youtube.com/watch?v=_eZNTBJ-Hyc', 'Valentina Monetta', 'Maybe'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2014), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 26, 'https://www.youtube.com/watch?v=C6w7tmWvm6M', 'Molly', 'Children of the Universe');

-- 2015 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2015,
    '2015-05-23 21:00:00+02',
    '2015-05-24 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Словения'), 1, 'https://www.youtube.com/watch?v=y1c7gDOaFY0', 'Maraaya', 'Here for You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Франция'), 2, 'https://www.youtube.com/watch?v=tJful-Jtc9o', 'Lisa Angell', 'N''oubliez pas'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Израиль'), 3, 'https://www.youtube.com/watch?v=-w-2K9U8UIs', 'Nadav Guedj', 'Golden Boy'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Эстония'), 4, 'https://www.youtube.com/watch?v=Xsx_dIJOtgI', 'Elina Born and Stig Rästa', 'Goodbye to Yesterday'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 5, 'https://www.youtube.com/watch?v=Ethyx4A5ajM', 'Electro Velvet', 'Still in Love with You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Армения'), 6, 'https://www.youtube.com/watch?v=Z6O8pr7HH94', 'Genealogy', 'Face the Shadow'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Литва'), 7, 'https://www.youtube.com/watch?v=xtB_slM63JA', 'Monika Linkytė and Vaidas Baumila', 'This Time'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 8, 'https://www.youtube.com/watch?v=gXGo70i94S8', 'Bojana Stamenov', 'Beauty Never Lies'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 9, 'https://www.youtube.com/watch?v=vY1FCXDojMU', 'Mørland and Debrah Scarlett', 'A Monster Like Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 10, 'https://www.youtube.com/watch?v=5sGOwFVUU0I', 'Måns Zelmerlöw', 'Heroes'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 11, 'https://www.youtube.com/watch?v=RUNEoMBLGc8', 'John Karayiannis', 'One Thing I Should Have Done'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Австралия'), 12, 'https://www.youtube.com/watch?v=H0EhhZWXTng', 'Guy Sebastian', 'Tonight Again'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Бельгия'), 13, 'https://www.youtube.com/watch?v=G48p8eNhnd8', 'Loïc Nottet', 'Rhythm Inside'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Австрия'), 14, 'https://www.youtube.com/watch?v=-f993p0CAV8', 'The Makemakes', 'I Am Yours'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Греция'), 15, 'https://www.youtube.com/watch?v=1RZ2JGC9Nlw', 'Maria Elena Kyriakou', 'One Last Breath'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Черногория'), 16, 'https://www.youtube.com/watch?v=MTyi2XeeoJU', 'Knez', 'Adio'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Германия'), 17, 'https://www.youtube.com/watch?v=DGMJOchTRPc', 'Ann Sophie', 'Black Smoke'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Польша'), 18, 'https://www.youtube.com/watch?v=TReYIZYloSg', 'Monika Kuszyńska', 'In the Name of Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Латвия'), 19, 'https://www.youtube.com/watch?v=-usdXbeGHi8', 'Aminata', 'Love Injected'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Румыния'), 20, 'https://www.youtube.com/watch?v=BGw0tUKOW7o', 'Voltaj', 'De la capăt'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Испания'), 21, 'https://www.youtube.com/watch?v=NsMiJJlElAY', 'Edurne', 'Amanecer'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 22, 'https://www.youtube.com/watch?v=MdybVsBESQc', 'Boggie', 'Wars for Nothing'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Грузия'), 23, 'https://www.youtube.com/watch?v=ddC8lviTERI', 'Nina Sublatti', 'Warrior'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 24, 'https://www.youtube.com/watch?v=IcVrAFtc5YI', 'Elnur Hüseynov', 'Hour of the Wolf'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Россия'), 25, 'https://www.youtube.com/watch?v=GVJW9ImpiWc', 'Polina Gagarina', 'A Million Voices'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Албания'), 26, 'https://www.youtube.com/watch?v=fsH9iVaZSaA', 'Elhaida Dani', 'I''m Alive'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2015), (SELECT id FROM countries WHERE name_ru = 'Италия'), 27, 'https://www.youtube.com/watch?v=1TOMqZV2jA8', 'Il Volo', 'Grande amore');

-- 2016 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2016,
    '2016-05-14 21:00:00+02',
    '2016-05-15 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Бельгия'), 1, 'https://www.youtube.com/watch?v=xDBjhAGaeWg', 'Laura Tesoro', 'What''s the Pressure'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Чехия'), 2, 'https://www.youtube.com/watch?v=wVy3rTNSN4Y', 'Gabriela Gunčíková', 'I Stand'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Нидерланды'), 3, 'https://www.youtube.com/watch?v=PwF9DI89q1w', 'Douwe Bob', 'Slow Down'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 4, 'https://www.youtube.com/watch?v=cB8eSLyffsI', 'Samra', 'Miracle'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 5, 'https://www.youtube.com/watch?v=bT2DcpqvorI', 'Freddie', 'Pioneer'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Италия'), 6, 'https://www.youtube.com/watch?v=YfnhDfy8AjU', 'Francesca Michielin', 'No Degree of Separation'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Израиль'), 7, 'https://www.youtube.com/watch?v=SS5TB2XUdgs', 'Hovi Star', 'Made of Stars'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Болгария'), 8, 'https://www.youtube.com/watch?v=PQqUTigWKHY', 'Poli Genova', 'If Love Was a Crime'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 9, 'https://www.youtube.com/watch?v=2AIPrVchJN8', 'Frans', 'If I Were Sorry'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Германия'), 10, 'https://www.youtube.com/watch?v=f-Z7pKopP9s', 'Jamie-Lee', 'Ghost'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Франция'), 11, 'https://www.youtube.com/watch?v=-aLPsiyavcU', 'Amir', 'J''ai cherché'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Польша'), 12, 'https://www.youtube.com/watch?v=87BBmxm7IJU', 'Michał Szpak', 'Color of Your Life'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Австралия'), 13, 'https://www.youtube.com/watch?v=5ymFX91HwM0', 'Dami Im', 'Sound of Silence'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 14, 'https://www.youtube.com/watch?v=idUxTqmCFR0', 'Minus One', 'Alter Ego'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 15, 'https://www.youtube.com/watch?v=yTlxLiSQWjk', 'Sanja Vučić Zaa', 'Goodbye (Shelter)'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Литва'), 16, 'https://www.youtube.com/watch?v=O1wS7NNCo38', 'Donny Montell', 'I''ve Been Waiting for This Night'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Хорватия'), 17, 'https://www.youtube.com/watch?v=8QUM-_EbE2o', 'Nina Kraljić', 'Lighthouse'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Россия'), 18, 'https://www.youtube.com/watch?v=e94dst20C9Y', 'Sergey Lazarev', 'You Are the Only One'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Испания'), 19, 'https://www.youtube.com/watch?v=EgbjNsxehrY', 'Barei', 'Say Yay!'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Латвия'), 20, 'https://www.youtube.com/watch?v=o-2Lt7zzlBU', 'Justs', 'Heartbeat'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Украина'), 21, 'https://www.youtube.com/watch?v=B-rnM-MwRHY', 'Jamala', '1944'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Мальта'), 22, 'https://www.youtube.com/watch?v=58YS5m3EYVY', 'Ira Losco', 'Walk on Water'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Грузия'), 23, 'https://www.youtube.com/watch?v=rviE2-9eiTI', 'Nika Kocharov and Young Georgian Lolitaz', 'Midnight Gold'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Австрия'), 24, 'https://www.youtube.com/watch?v=3xuPYt5i5cE', 'Zoë', 'Loin d''ici'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 25, 'https://www.youtube.com/watch?v=XgGahKsBtts', 'Joe and Jake', 'You''re Not Alone'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2016), (SELECT id FROM countries WHERE name_ru = 'Армения'), 26, 'https://www.youtube.com/watch?v=nL066Rp7J7k', 'Iveta Mukuchyan', 'LoveWave');

-- 2017 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2017,
    '2017-05-13 21:00:00+02',
    '2017-05-14 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Израиль'), 1, 'https://www.youtube.com/watch?v=YapnW_zfWJk', 'Imri', 'I Feel Alive'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Польша'), 2, 'https://www.youtube.com/watch?v=xVCokUW3BUo', 'Kasia Moś', 'Flashlight'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Беларусь'), 3, 'https://www.youtube.com/watch?v=coFfFZvPmRY', 'Naviband', 'Story of My Life'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Австрия'), 4, 'https://www.youtube.com/watch?v=gJIjZF5sgzg', 'Nathan Trent', 'Running on Air'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Армения'), 5, 'https://www.youtube.com/watch?v=g373Libigw8', 'Artsvik', 'Fly with Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Нидерланды'), 6, 'https://www.youtube.com/watch?v=_9XsCqGg8ls', 'Og3ne', 'Lights and Shadows'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Молдова'), 7, 'https://www.youtube.com/watch?v=SWaQdHoCvYk', 'SunStroke Project', 'Hey Mamma'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 8, 'https://www.youtube.com/watch?v=hqM0AkP7zcI', 'Joci Pápai', 'Origo'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Италия'), 9, 'https://www.youtube.com/watch?v=KieE_MLv-ZY', 'Francesco Gabbani', 'Occidentali''s Karma'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Дания'), 10, 'https://www.youtube.com/watch?v=vUbGnq8maS0', 'Anja', 'Where I Am'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Португалия'), 11, 'https://www.youtube.com/watch?v=Qotooj7ODCM', 'Salvador Sobral', 'Amar pelos dois'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 12, 'https://www.youtube.com/watch?v=g_HSRHqd_7s', 'Dihaj', 'Skeletons'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Хорватия'), 13, 'https://www.youtube.com/watch?v=4Sya_GIxsw0', 'Jacques Houdek', 'My Friend'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Австралия'), 14, 'https://www.youtube.com/watch?v=YTCP4S1dqwA', 'Isaiah', 'Don''t Come Easy'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Греция'), 15, 'https://www.youtube.com/watch?v=iO14p2uCzkM', 'Demy', 'This Is Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Испания'), 16, 'https://www.youtube.com/watch?v=qAOXHdLdlqQ', 'Manel Navarro', 'Do It for Your Lover'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 17, 'https://www.youtube.com/watch?v=U-Bk0PCotY0', 'Jowst', 'Grab the Moment'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 18, 'https://www.youtube.com/watch?v=Fl1GYTg4GmA', 'Lucie Jones', 'Never Give Up on You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 19, 'https://www.youtube.com/watch?v=ZDe9FCLgE04', 'Hovig', 'Gravity'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Румыния'), 20, 'https://www.youtube.com/watch?v=3lop3KkfRPA', 'Ilinca feat. Alex Florea', 'Yodel It!'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Германия'), 21, 'https://www.youtube.com/watch?v=t1KKjykH6fw', 'Levina', 'Perfect Life'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Украина'), 22, 'https://www.youtube.com/watch?v=mogNLkIEI0k', 'O.Torvald', 'Time'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Бельгия'), 23, 'https://www.youtube.com/watch?v=oxsCmChDYwA', 'Blanche', 'City Lights'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 24, 'https://www.youtube.com/watch?v=yqSV81kmsbQ', 'Robin Bengtsson', 'I Can''t Go On'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Болгария'), 25, 'https://www.youtube.com/watch?v=OMmm-G078LM', 'Kristian Kostov', 'Beautiful Mess'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2017), (SELECT id FROM countries WHERE name_ru = 'Франция'), 26, 'https://www.youtube.com/watch?v=HjiLKsLQl4I', 'Alma', 'Requiem');

-- 2018 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2018,
    '2018-05-12 21:00:00+02',
    '2018-05-13 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Украина'), 1, 'https://www.youtube.com/watch?v=lhOuMYtQ94M', 'Mélovin', 'Under the Ladder'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Испания'), 2, 'https://www.youtube.com/watch?v=28oMEFgSoyk', 'Amaia and Alfred', 'Tu canción'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Словения'), 3, 'https://www.youtube.com/watch?v=fBU0_hfST_o', 'Lea Sirk', 'Hvala, ne!'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Литва'), 4, 'https://www.youtube.com/watch?v=pPDFhfU25go', 'Ieva Zasimauskaitė', 'When We''re Old'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Австрия'), 5, 'https://www.youtube.com/watch?v=a8Yvzo1puoE', 'Cesár Sampson', 'Nobody but You'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Эстония'), 6, 'https://www.youtube.com/watch?v=ImawXdXIGd8', 'Elina Nechayeva', 'La forza'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 7, 'https://www.youtube.com/watch?v=N17BdRwOvBA', 'Alexander Rybak', 'That''s How You Write a Song'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Португалия'), 8, 'https://www.youtube.com/watch?v=mHjGbADMcOI', 'Cláudia Pascoal', 'O jardim'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 9, 'https://www.youtube.com/watch?v=n0ydtXG-tqY', 'SuRie', 'Storm'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 10, 'https://www.youtube.com/watch?v=WRQWQShbPHU', 'Sanja Ilić and Balkanika', 'Nova deca'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Германия'), 11, 'https://www.youtube.com/watch?v=N4YwSy4Q9k0', 'Michael Schulte', 'You Let Me Walk Alone'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Албания'), 12, 'https://www.youtube.com/watch?v=uyl1b1fVmdU', 'Eugent Bushpepa', 'Mall'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Франция'), 13, 'https://www.youtube.com/watch?v=vjLuTtUv0Ns', 'Madame Monsieur', 'Mercy'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Чехия'), 14, 'https://www.youtube.com/watch?v=w0ZUuQe0shk', 'Mikolas Josef', 'Lie to Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Дания'), 15, 'https://www.youtube.com/watch?v=c0kP074Mw0k', 'Rasmussen', 'Higher Ground'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Австралия'), 16, 'https://www.youtube.com/watch?v=2BBLTY6T2Ck', 'Jessica Mauboy', 'We Got Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Финляндия'), 17, 'https://www.youtube.com/watch?v=hdjiuCUUpPc', 'Saara Aalto', 'Monsters'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Болгария'), 18, 'https://www.youtube.com/watch?v=vC-1tdGDFQc', 'Equinox', 'Bones'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Молдова'), 19, 'https://www.youtube.com/watch?v=5evrVL1rJPA', 'DoReDoS', 'My Lucky Day'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 20, 'https://www.youtube.com/watch?v=z101xtUBflU', 'Benjamin Ingrosso', 'Dance You Off'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Венгрия'), 21, 'https://www.youtube.com/watch?v=p6e1TmYb33w', 'AWS', 'Viszlát nyár'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Израиль'), 22, 'https://www.youtube.com/watch?v=84LBjXaeKk4', 'Netta', 'Toy'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Нидерланды'), 23, 'https://www.youtube.com/watch?v=p7jHz5asMs8', 'Waylon', 'Outlaw in ''Em'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Ирландия'), 24, 'https://www.youtube.com/watch?v=QD1cQ2wZ1dk', 'Ryan O''Shaughnessy', 'Together'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 25, 'https://www.youtube.com/watch?v=vyDTbJ4wenY', 'Eleni Foureira', 'Fuego'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2018), (SELECT id FROM countries WHERE name_ru = 'Италия'), 26, 'https://www.youtube.com/watch?v=81M-mp5t8uM', 'Ermal Meta and Fabrizio Moro', 'Non mi avete fatto niente');

-- 2019 Grand Final
INSERT INTO contests (id, type, year, starts, ends)
VALUES (
    gen_random_uuid(),
    'final',
    2019,
    '2019-05-18 21:00:00+02',
    '2019-05-19 01:00:00+02'
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
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Мальта'), 1, 'https://www.youtube.com/watch?v=H-byNyainRA', 'Michela', 'Chameleon'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Албания'), 2, 'https://www.youtube.com/watch?v=-NAbYUoxIfg', 'Jonida Maliqi', 'Ktheju tokës'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Чехия'), 3, 'https://www.youtube.com/watch?v=NJxLKFZKhNs', 'Lake Malawi', 'Friend of a Friend'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Германия'), 4, 'https://www.youtube.com/watch?v=yFbSj_yPCPM', 'Sisters', 'Sister'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Россия'), 5, 'https://www.youtube.com/watch?v=eNzlxEZ_JG4', 'Sergey Lazarev', 'Scream'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Дания'), 6, 'https://www.youtube.com/watch?v=mucmvTLHJQM', 'Leonora', 'Love Is Forever'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Сан-Марино'), 7, 'https://www.youtube.com/watch?v=d_iEis_ZRoQ', 'Serhat', 'Say Na Na Na'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Северная Македония'), 8, 'https://www.youtube.com/watch?v=0pAURAsum4A', 'Tamara Todevska', 'Proud'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Швеция'), 9, 'https://www.youtube.com/watch?v=oEdWkdVKIqQ', 'John Lundvik', 'Too Late for Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Словения'), 10, 'https://www.youtube.com/watch?v=vGtPNQ6g6Ng', 'Zala Kralj and Gašper Šantl', 'Sebi'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Кипр'), 11, 'https://www.youtube.com/watch?v=wIpEI-NhRxg', 'Tamta', 'Replay'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Нидерланды'), 12, 'https://www.youtube.com/watch?v=R3D-r4ogr7s', 'Duncan Laurence', 'Arcade'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Греция'), 13, 'https://www.youtube.com/watch?v=5pQV_-40m_U', 'Katerine Duska', 'Better Love'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Израиль'), 14, 'https://www.youtube.com/watch?v=uCdS42uCqP0', 'Kobi Marimi', 'Home'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Норвегия'), 15, 'https://www.youtube.com/watch?v=Ovt7YGHAj8I', 'Keiino', 'Spirit in the Sky'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Великобритания'), 16, 'https://www.youtube.com/watch?v=HV-eOhTS8Dw', 'Michael Rice', 'Bigger than Us'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Исландия'), 17, 'https://www.youtube.com/watch?v=kTb69WkBbvs', 'Hatari', 'Hatrið mun sigra'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Эстония'), 18, 'https://www.youtube.com/watch?v=PppvSl3_W4E', 'Victor Crone', 'Storm'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Беларусь'), 19, 'https://www.youtube.com/watch?v=C1Y9LBusHzA', 'Zena', 'Like It'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Азербайджан'), 20, 'https://www.youtube.com/watch?v=D6vvuCiHwSs', 'Chingiz', 'Truth'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Франция'), 21, 'https://www.youtube.com/watch?v=AjZNbIQtWvw', 'Bilal Hassani', 'Roi'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Италия'), 22, 'https://www.youtube.com/watch?v=M-aoyPa41Ic', 'Mahmood', 'Soldi'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Сербия'), 23, 'https://www.youtube.com/watch?v=X0mdaeAyPvA', 'Nevena Božović', 'Kruna'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Швейцария'), 24, 'https://www.youtube.com/watch?v=8aUYzqAIdoM', 'Luca Hänni', 'She Got Me'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Австралия'), 25, 'https://www.youtube.com/watch?v=Eo5H62mCIsg', 'Kate Miller-Heidke', 'Zero Gravity'),
(gen_random_uuid(), (SELECT id FROM contests WHERE type = 'final' AND year = 2019), (SELECT id FROM countries WHERE name_ru = 'Испания'), 26, 'https://www.youtube.com/watch?v=lWXuohnM14U', 'Miki', 'La venda');
