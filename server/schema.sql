CREATE TABLE accounts (
  uuid TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  hashed_password TEXT NOT NULL
);