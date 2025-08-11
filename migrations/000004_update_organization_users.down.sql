ALTER TABLE organization_users
DROP COLUMN role;

DROP TYPE IF EXISTS role_enum;