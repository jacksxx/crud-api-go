CREATE SCHEMA IF NOT EXISTS prod;

SET
    timezone = 'America/Bahia';

CREATE TABLE
    IF NOT EXISTS prod.categorias (
        categorias_id SERIAL PRIMARY KEY,
        categorias_name VARCHAR(200) NOT NULL
    );

-- Inserindo as categorias
INSERT INTO
    prod.categorias (categorias_name)
VALUES
    ('Bebidas'),
    ('Lanches'),
    ('Doces e Sobremesas'),
    ('Produtos de Padaria'),
    ('Alimentos Básicos'),
    ('Carnes e Proteínas'),
    ('Frutas e Verduras'),
    ('Produtos de Limpeza');

CREATE TABLE
    IF NOT EXISTS prod.products (
        products_id SERIAL PRIMARY KEY,
        products_name VARCHAR(50) NOT NULL,
        products_price DECIMAL(10, 2) NOT NULL,
        data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        categorias_id INT NOT NULL,
        categorias_name VARCHAR(200) NOT NULL,
        FOREIGN KEY (categorias_id) REFERENCES prod.categorias (categorias_id) ON DELETE CASCADE
    );

-- Inserindo os produtos com categorias
INSERT INTO
    prod.products (
        products_name,
        products_price,
        categorias_id,
        categorias_name
    )
VALUES
    -- Bebidas (categoria_id = 1)
    ('Coca-Cola', 4.50, 1, 'Bebidas'),
    ('Suco de Laranja', 3.20, 1, 'Bebidas'),
    ('Água Mineral', 1.00, 1, 'Bebidas'),
    ('Cerveja Heineken', 6.80, 1, 'Bebidas'),
    ('Refrigerante Guaraná', 3.90, 1, 'Bebidas'),
    ('Chá Gelado', 2.80, 1, 'Bebidas'),
    -- Lanches (categoria_id = 2)
    ('X-Tudo', 8.90, 2, 'Lanches'),
    ('Pastel de Queijo', 4.50, 2, 'Lanches'),
    ('Coxinha', 3.00, 2, 'Lanches'),
    ('Hambúrguer', 7.99, 2, 'Lanches'),
    ('Misto Quente', 5.00, 2, 'Lanches'),
    ('Esfiha', 4.25, 2, 'Lanches'),
    -- Doces e Sobremesas (categoria_id = 3)
    (
        'Chocolate ao Leite',
        7.50,
        3,
        'Doces e Sobremesas'
    ),
    (
        'Sorvete de Morango',
        5.25,
        3,
        'Doces e Sobremesas'
    ),
    ('Pudim', 4.75, 3, 'Doces e Sobremesas'),
    ('Brigadeiro', 2.50, 3, 'Doces e Sobremesas'),
    ('Torta de Limão', 6.90, 3, 'Doces e Sobremesas'),
    (
        'Bolo de Chocolate',
        9.99,
        3,
        'Doces e Sobremesas'
    ),
    -- Produtos de Padaria (categoria_id = 4)
    ('Pão Francês', 0.80, 4, 'Produtos de Padaria'),
    ('Croissant', 3.99, 4, 'Produtos de Padaria'),
    (
        'Bolo de Cenoura',
        10.50,
        4,
        'Produtos de Padaria'
    ),
    ('Rosquinha', 2.75, 4, 'Produtos de Padaria'),
    ('Pão de Queijo', 4.50, 4, 'Produtos de Padaria'),
    ('Baguete', 5.20, 4, 'Produtos de Padaria'),
    -- Alimentos Básicos (categoria_id = 5)
    ('Arroz 5kg', 20.00, 5, 'Alimentos Básicos'),
    ('Feijão 1kg', 9.50, 5, 'Alimentos Básicos'),
    (
        'Macarrão Espaguete',
        3.25,
        5,
        'Alimentos Básicos'
    ),
    ('Óleo de Soja', 7.99, 5, 'Alimentos Básicos'),
    ('Sal', 2.30, 5, 'Alimentos Básicos'),
    ('Açúcar', 4.00, 5, 'Alimentos Básicos'),
    -- Carnes e Proteínas (categoria_id = 6)
    (
        'Filé de Frango 1kg',
        18.50,
        6,
        'Carnes e Proteínas'
    ),
    ('Carne Moída 1kg', 27.90, 6, 'Carnes e Proteínas'),
    ('Salmão 500g', 32.80, 6, 'Carnes e Proteínas'),
    (
        'Linguiça Toscana',
        15.99,
        6,
        'Carnes e Proteínas'
    ),
    (
        'Bife de Contra Filé',
        29.50,
        6,
        'Carnes e Proteínas'
    ),
    ('Peito de Peru', 25.75, 6, 'Carnes e Proteínas'),
    -- Frutas e Verduras (categoria_id = 7)
    ('Maçã', 4.20, 7, 'Frutas e Verduras'),
    ('Tomate', 5.30, 7, 'Frutas e Verduras'),
    ('Alface', 2.50, 7, 'Frutas e Verduras'),
    ('Banana', 3.00, 7, 'Frutas e Verduras'),
    ('Abacaxi', 7.20, 7, 'Frutas e Verduras'),
    ('Cenoura', 4.80, 7, 'Frutas e Verduras'),
    -- Produtos de Limpeza (categoria_id = 8)
    ('Detergente', 3.60, 8, 'Produtos de Limpeza'),
    ('Sabão em Pó', 15.90, 8, 'Produtos de Limpeza'),
    (
        'Papel Higiênico 12un',
        18.00,
        8,
        'Produtos de Limpeza'
    ),
    ('Amaciante', 12.50, 8, 'Produtos de Limpeza'),
    ('Álcool 70%', 8.99, 8, 'Produtos de Limpeza'),
    (
        'Esponja de Limpeza',
        2.75,
        8,
        'Produtos de Limpeza'
    );