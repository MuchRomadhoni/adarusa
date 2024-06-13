CREATE TABLE IF NOT EXISTS pinjaman(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    admin_id int,
    jumlah_pinjaman float NOT NULL,
    tenor int NOT NULL,
    bunga_per_bulan float,
    description TEXT,
    status_pengajuan pinjaman_status DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_admin_id FOREIGN KEY (admin_id) REFERENCES users(id)
)