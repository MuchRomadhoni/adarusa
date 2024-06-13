CREATE TYPE withdrawal_status AS ENUM('pending','paid','rejected');

CREATE TABLE IF NOT EXISTS withdrawal(
  id SERIAL PRIMARY KEY,
  user_id int NOT NULL,
  amount float NOT NULL,
  status withdrawal_status DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
)