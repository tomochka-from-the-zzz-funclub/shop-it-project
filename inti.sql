-- Расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Пользователи
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Информация о покупателях
CREATE TABLE buyer_info (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number BIGINT NOT NULL,
    gender BOOLEAN NOT NULL,
    birthdate TIMESTAMP NOT NULL,
    CONSTRAINT fk_buyer_user FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);

-- Продавцы
CREATE TABLE sellers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Информация о продавцах
CREATE TABLE seller_info (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number BIGINT NOT NULL,
    gender BOOLEAN NOT NULL,
    birthdate TIMESTAMP NOT NULL,
    CONSTRAINT fk_seller_user FOREIGN KEY (id) REFERENCES sellers(id) ON DELETE CASCADE
);

-- Карточки товаров
CREATE TABLE good_cards (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price NUMERIC(10, 2) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    weight NUMERIC(10, 2) NOT NULL,
    seller_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT fk_good_seller FOREIGN KEY (seller_id) REFERENCES sellers(id) ON DELETE CASCADE
);

-- Склад товаров
CREATE TABLE goods (
    card_id UUID UNIQUE REFERENCES good_cards(uuid) ON DELETE CASCADE,
    quantity INT DEFAULT 0
);

-- Таблица заказов
CREATE TABLE orders (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount NUMERIC(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Created',
    CONSTRAINT fk_order_customer FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (status IN ('Created', 'ReadyForPickup', 'Received'))
);

-- Позиции в заказе
CREATE TABLE order_items (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    good_uuid UUID NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(uuid) ON DELETE CASCADE,
    FOREIGN KEY (good_uuid) REFERENCES goods(card_id) ON DELETE CASCADE
);

-- Корзина
CREATE TABLE bag (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    goods UUID[] NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE
);
