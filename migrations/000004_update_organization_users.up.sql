CREATE TYPE role_enum AS ENUM (
	'admin',
	'member'
);

ALTER TABLE organization_users
ADD COLUMN role role_enum NOT NULL DEFAULT 'member';