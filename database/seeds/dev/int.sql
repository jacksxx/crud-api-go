CREATE SCHEMA IF NOT EXISTS prod;

CREATE TABLE
    IF NOT EXISTS prod.products (
        products_id SERIAL PRIMARY KEY,
        products_name VARCHAR(50) NOT NULL,
        products_price DECIMAL(10, 2) NOT NULL
    );

INSERT INTO
    prod.products (products_name, products_price)
VALUES
    ('Biscoito', 3.50),
    ('Banana', 2.00),
    ('BK Whopper', 6.90),
    ('MC Big Mac', 5.99);