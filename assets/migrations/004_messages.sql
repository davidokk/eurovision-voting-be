create table messages (
	contest_id uuid REFERENCES contests(id),
	user_id uuid REFERENCES users(id),
	message text,
	created_at TIMESTAMP WITH TIME ZONE
);

alter table messages add column type text;
alter table messages add column performance_id uuid;

alter table messages add column score int;
alter table messages add column old_score int;
alter table messages add column comment text;
alter table messages add column gif text;
