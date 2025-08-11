CREATE TABLE IF NOT EXISTS organization_users (
	ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id UUID NOT NULL,
	user_id UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	UNIQUE (organization_id, user_id),

	FOREIGN KEY (organization_id) 
		REFERENCES organizations(id) 
		ON DELETE CASCADE,
	FOREIGN KEY (user_id) 
		REFERENCES users(id) 
		ON DELETE CASCADE
);