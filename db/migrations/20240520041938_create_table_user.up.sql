CREATE TABLE IF NOT EXISTS "users" (
    id SERIAL PRIMARY KEY,
    email varchar(50) NOT NULL,
    password varchar(50) NOT NULL,
    role_id INT NOT NULL,
    status user_status DEFAULT 'unverified',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES roles(id)
)