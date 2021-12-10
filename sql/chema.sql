-- Товары
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    qty INTEGER NOT NULL DEFAULT 0 CHECK (qty >= 0),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 
-- 
-- Сотрудники
CREATE TABLE managers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    salary INTEGER NOT NULL CHECK (salary > 0),
    plan INTEGER NOT NULL CHECK (plan >= 0),
    boss_id BIGINT REFERENCES managers,
    department TEXT,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 
-- 
-- Покупатели
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 
-- 
CREATE TABLE customers_tokens (
    token TEXT NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL REFERENCES customers,
    expire TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 
--  
-- Продажи
CREATE TABLE sales (
    id BIGSERIAL PRIMARY KEY,
    manager_id BIGINT NOT NULL REFERENCES managers,
    customer_id BIGINT REFERENCES customers,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 
-- 
-- Позиции
CREATE TABLE sale_positions (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES sales,
    product_id BIGINT NOT NULL REFERENCES products,
    name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    qty INTEGER NOT NULL CHECK (qty > 0),
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 
-- 
DROP TABLE sale_positions;
DROP TABLE sales;
DROP TABLE products;
DROP TABLE managers;
DROP TABLE customers;
DROP TABLE customers_tokens;
-- 
-- 
UPDATE customers_tokens
SET expire = '2021-12-10 10:11:16.651102';