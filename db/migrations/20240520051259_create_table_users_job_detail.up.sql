CREATE TABLE IF NOT EXISTS users_job_detail(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    job_name varchar(50) NOT NULL,
    gaji float(50) NOT NULL,
    office_name varchar(50) NOT NULL,
    office_contact varchar(50) NOT NULL,
    address text not null,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
)