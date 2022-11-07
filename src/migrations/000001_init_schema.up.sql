DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id serial PRIMARY KEY,
  email VARCHAR(128) UNIQUE NOT NULL,
  username VARCHAR(128) UNIQUE NOT NULL,
  password_hash VARCHAR(128) NOT NULL,
  first_name VARCHAR (64),
  last_name VARCHAR (64),
  token VARCHAR (128),
  status SMALLINT NOT NULL,
  token_expires_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

/* INITIALIZE USERS */
INSERT INTO users (id, email, username, password_hash, first_name, last_name, token, status, token_expires_at) VALUES
  (1, 'admin@kallabs.by', 'admin', '$2a$14$1uz8bdnCERhrFJ1qDZ0gwOxxmHy4NuYsAu2mckpzL3r5C7WbO3nCO', '', '', '', 2, NOW());
