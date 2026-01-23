-- +goose Up
CREATE TABLE links(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  code TEXT NOT NULL UNIQUE,
  link_url TEXT NOT NULL,
  user_id UUID NOT NULL,
  CONSTRAINT fk_users
  FOREIGN KEY (user_id)
  REFERENCES users(id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE links;
