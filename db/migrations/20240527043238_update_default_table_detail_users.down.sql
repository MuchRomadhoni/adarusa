ALTER TABLE detail_users
    ALTER COLUMN limit_id DROP NOT NULL,
    ALTER COLUMN foto_ktp DROP NOT NULL,
    ALTER COLUMN foto_selfie DROP NOT NULL;