CREATE TABLE IF NOT EXISTS claim_tugas(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    collector_id INT NOT NULL,
    status claim_status,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_collector_id FOREIGN KEY (collector_id) REFERENCES users(id)
)