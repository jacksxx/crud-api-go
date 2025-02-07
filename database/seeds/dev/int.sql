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
    -- Bebidas
    ('Coca-Cola', 4.50),
    ('Suco de Laranja', 3.20),
    ('Água Mineral', 1.00),
    ('Cerveja Heineken', 6.80),
    -- Lanches
    ('X-Tudo', 8.90),
    ('Pastel de Queijo', 4.50),
    ('Coxinha', 3.00),
    -- Doces e Sobremesas
    ('Chocolate ao Leite', 7.50),
    ('Sorvete de Morango', 5.25),
    ('Pudim', 4.75),
    -- Produtos de Padaria
    ('Pão Francês', 0.80),
    ('Croissant', 3.99),
    ('Bolo de Cenoura', 10.50),
    -- Alimentos Básicos
    ('Arroz 5kg', 20.00),
    ('Feijão 1kg', 9.50),
    ('Macarrão Espaguete', 3.25),
    ('Óleo de Soja', 7.99),
    -- Carnes e Proteínas
    ('Filé de Frango 1kg', 18.50),
    ('Carne Moída 1kg', 27.90),
    ('Salmão 500g', 32.80),
    -- Frutas e Verduras
    ('Maçã', 4.20),
    ('Tomate', 5.30),
    ('Alface', 2.50),
    -- Produtos Industriais
    ('Detergente', 3.60),
    ('Sabão em Pó', 15.90),
    ('Papel Higiênico 12un', 18.00);