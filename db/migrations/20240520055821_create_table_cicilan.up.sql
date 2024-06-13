CREATE TABLE IF NOT EXISTS cicilan(
    id SERIAL PRIMARY KEY,
    pinjaman_id int NOT NULL,
    tanggal_jatuh_tempo TIMESTAMP NOT NULL,
    tanggal_selesai_bayar TIMESTAMP,
    jumlah_bayar float NOT NULL,
    status varchar(50) NOT NULL,
    CONSTRAINT fk_pinjaman_id FOREIGN KEY (pinjaman_id) REFERENCES pinjaman(id)
)