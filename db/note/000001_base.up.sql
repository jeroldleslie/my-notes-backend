CREATE DATABASE notes;
\connect notes;


CREATE EXTENSION pg_trgm;

-- tables


-- Drop table

-- DROP TABLE public.users

CREATE TABLE users (
	id bigserial NOT NULL,
	"name" text NULL,
	email text NOT NULL,
	"password" text NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (
	OIDS=FALSE
) ;

-- DROP TABLE public.notes
CREATE TABLE notes (
	id bigserial NOT NULL,
	title text NULL,
	"content" text NULL,
	priority text NULL,
	created_at timestamp NULL DEFAULT now(),
	updated_at timestamp NULL DEFAULT now(),
	remind_from timestamp NULL,
	remind_until timestamp NULL,
	image bytea NULL,
	user_id int8 NOT NULL,
	color text NULL DEFAULT '#ffffff'::text,
	CONSTRAINT notes_pkey PRIMARY KEY (id),
	CONSTRAINT notes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
)
WITH (
	OIDS=FALSE
) ;
CREATE INDEX notes_content_gin_trgm_idx ON public.notes USING gin (content gin_trgm_ops) ;
CREATE INDEX notes_title_gin_trgm_idx ON public.notes USING gin (title gin_trgm_ops) ;


-- Drop table

-- DROP TABLE public.files

CREATE TABLE files (
	id bigserial NOT NULL,
	file_name text NULL,
	"content" text NULL,
	content_type text NULL,
	created_at timestamp NULL DEFAULT now(),
	updated_at timestamp NULL DEFAULT now(),
	note_id int8 NOT NULL,
	CONSTRAINT files_note_id_key UNIQUE (note_id),
	CONSTRAINT files_pkey PRIMARY KEY (id),
	CONSTRAINT files_note_id_fkey FOREIGN KEY (note_id) REFERENCES notes(id)
)
WITH (
	OIDS=FALSE
) ;





