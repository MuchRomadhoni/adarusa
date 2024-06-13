CREATE TABLE IF NOT EXISTS kontak_darurat(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    name varchar(50) NOT NULL,
    phone_number varchar(50) NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
)