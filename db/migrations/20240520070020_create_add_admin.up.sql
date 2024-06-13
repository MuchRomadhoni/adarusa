CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE users ALTER COLUMN password TYPE VARCHAR(255);

INSERT INTO users(email,password,role_id,status) VALUES ('admin@mail.com', crypt('admin123', gen_salt('bf')),1,'verified')