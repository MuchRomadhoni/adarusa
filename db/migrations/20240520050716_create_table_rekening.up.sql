CREATE TABLE IF NOT EXISTS rekening(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    account_number varchar(50) NOT NULL,
    bank_name varchar(50) NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
)