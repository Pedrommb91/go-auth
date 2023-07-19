-- +goose Up
-- +goose StatementBegin
CREATE TABLE credentials (
  id SERIAL PRIMARY KEY, 
  salt VARCHAR(254) NOT NULL,
  passhash VARCHAR(254) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL 
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(63) UNIQUE NOT NULL,
  email VARCHAR(254) UNIQUE NOT NULL
    CONSTRAINT 
      proper_email CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
  credentials_id INT
    CONSTRAINT fk_users_credentials 
      REFERENCES credentials
      ON UPDATE CASCADE ON DELETE CASCADE,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
DROP TABLE credentials
-- +goose StatementEnd
