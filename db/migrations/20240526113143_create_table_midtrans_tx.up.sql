CREATE TABLE IF NOT EXISTS midtrans_tx(
    id SERIAL PRIMARY KEY,
    cicilan_id int NOT NULL,
    amount float NOT NULL,
    snap_url varchar(255) NOT NULL,
    status midtrans_status NOT NULL DEFAULT 'pending',
    CONSTRAINT fk_cicilan_id FOREIGN KEY (cicilan_id) REFERENCES cicilan(id)
)