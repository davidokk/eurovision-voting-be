create table messages (
	contest_id uuid REFERENCES contests(id),
	user_id uuid REFERENCES users(id),
	message text,
	created_at TIMESTAMP WITH TIME ZONE
);