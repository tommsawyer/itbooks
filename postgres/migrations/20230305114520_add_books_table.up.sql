CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  isbn TEXT UNIQUE NOT NULL,
  url TEXT NOT NULL,
  title TEXT NOT NULL,
  authors TEXT[],
  image TEXT NOT NULL,
  description TEXT NOT NULL,
  properties JSONB,
  publisher TEXT NOT NULL,
  published BOOLEAN NOT NULL DEFAULT FALSE,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW()
);
