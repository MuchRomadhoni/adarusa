CREATE TABLE IF NOT EXISTS balance(
    id SERIAL PRIMARY KEY,
    user_id int UNIQUE NOT NULL,
    amount float NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
)