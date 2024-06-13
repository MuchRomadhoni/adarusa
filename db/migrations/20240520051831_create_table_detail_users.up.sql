CREATE TABLE IF NOT EXISTS detail_users(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    limit_id int NOT NULL,
    NIK varchar(50) NOT NULL,
    fullname varchar(50) NOT NULL,
    phone_number varchar(50) NOT NULL,
    address text NOT NULL,
    city varchar(50) NOT NULL,
    foto_ktp varchar(50) NOT NULL,
    foto_selfie varchar(50) NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_limit_id FOREIGN KEY (limit_id) REFERENCES limit_pinjaman(id)
)