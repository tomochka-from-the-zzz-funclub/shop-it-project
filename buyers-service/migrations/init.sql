SET search_path TO public;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA public;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS buyer_info (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number BIGINT NOT NULL,
    gender BOOLEAN NOT NULL,
    birthdate TIMESTAMP NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

COMMENT ON TABLE users IS 'Таблица для хранения пользователей';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя (UUID)';
COMMENT ON COLUMN users.email IS 'Email пользователя (уникальный)';
COMMENT ON COLUMN users.password_hash IS 'Хэш пароля пользователя';
COMMENT ON COLUMN users.created_at IS 'Дата создания пользователя';

COMMENT ON TABLE buyer_info IS 'Таблица для хранения информации о покупателях';
COMMENT ON COLUMN buyer_info.id IS 'Идентификатор пользователя и покупателя (UUID, совпадает с users.id)';
COMMENT ON COLUMN buyer_info.name IS 'Имя покупателя';
COMMENT ON COLUMN buyer_info.phone_number IS 'Номер телефона покупателя';
COMMENT ON COLUMN buyer_info.gender IS 'Пол покупателя (true = женский, false = мужской)';
COMMENT ON COLUMN buyer_info.birthdate IS 'Дата рождения покупателя';

DO $$
BEGIN
    RAISE NOTICE 'Tables users and buyer_info created successfully';
END $$;