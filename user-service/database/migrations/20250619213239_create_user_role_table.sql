-- migrate:up
CREATE TABLE IF NOT EXISTS user_role (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);


-- migrate:down
DROP TABLE IF EXISTS "user_role";
