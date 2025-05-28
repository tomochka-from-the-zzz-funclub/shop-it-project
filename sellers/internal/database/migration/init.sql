CREATE TABLE sellers (
	uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    passport INT UNIQUE NOT NULL,
    telephone_number VARCHAR(16) UNIQUE NOT NULL,
    description TEXT
);
CREATE TABLE good_cards (
		uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Уникальный идентификатор карточки товара
		price NUMERIC(10, 2) NOT NULL,                  -- Цена товара
		name VARCHAR(255) NOT NULL,                     -- Название товара
		description TEXT NOT NULL,                      -- Описание товара
		weight NUMERIC(10, 2) NOT NULL,                -- Вес товара
		seller_id UUID NOT NULL,                        -- Уникальный идентификатор продавца
		is_active BOOLEAN NOT NULL DEFAULT TRUE         -- Статус активации товара
	);

	CREATE TABLE goods (
		card_id UUID REFERENCES good_cards(uuid) ON DELETE CASCADE, -- Внешний ключ на карточку товара
		quantity INT    DEFAULT 0                    -- Количество товара
	);