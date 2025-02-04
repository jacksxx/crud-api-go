CREATE TABLE
    products (
        products_id SERIAL PRIMARY KEY,
        products_name VARCHAR(50) NOT NULL,
        products_price DECIMAL(10, 2) NOT NULL
    )
INSERT INTO
    products (products_name, products_price)
VALUES
    ('Biscoito', 3.50),
    ('Banana', 2.00)