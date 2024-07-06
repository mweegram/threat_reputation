CREATE TABLE IF NOT EXISTS threats(
	id serial primary key,
	filename text not null,
	sha256 text not null,
	comments int[],
	submitted text not null
);

CREATE TABLE IF NOT EXISTS comments(
	id serial primary key,
	text text not null,
	date text not null
);