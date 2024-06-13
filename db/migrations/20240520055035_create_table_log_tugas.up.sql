CREATE TABLE IF NOT EXISTS log_tugas(
    id SERIAL PRIMARY KEY,
    tugas_id int NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_tugas_id FOREIGN KEY (tugas_id) REFERENCES claim_tugas(id)
)