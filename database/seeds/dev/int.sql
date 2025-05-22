CREATE SCHEMA IF NOT EXISTS prod;

SET
    timezone = 'America/Bahia';

CREATE TABLE 
    IF NOT EXISTS prod.usuarios (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    usuario VARCHAR(255) NOT NULL,
    email VARCHAR(250) NOT NULL,
    senha VARCHAR(100) NOT NULL,
    salt VARCHAR(100) NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT true,
    data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP NULL,
    data_inativacao TIMESTAMP NULL
);

INSERT INTO prod.usuarios (
    nome,
    usuario,
    email,
    senha,
    salt
) VALUES (
    'admin',
    'admin',
    'admin@example.com',
    '$2a$10$HnClX61FTm9i110IVXZubeAKD1ZcTHBdcrP/pgOrwT7RrXtQRzcIO',
    'YzkwVOZQhdLXRdX+I1k2PW0J13uIRTfKDbxbyhE7Xgk='
);

CREATE TABLE
    IF NOT EXISTS prod.categorias (
        categorias_id SERIAL PRIMARY KEY,
        categorias_name VARCHAR(200) NOT NULL,
        categorias_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        categorias_data_atualizacao TIMESTAMP NULL,
        categorias_data_inativacao TIMESTAMP NULL,
        categorias_status VARCHAR(50) NOT NULL DEFAULT 'ativo'
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
    IF NOT EXISTS prod.subcategorias (
        subcategorias_id SERIAL PRIMARY KEY,
        subcategorias_name VARCHAR(200) NOT NULL,
        categorias_id INT NOT NULL REFERENCES prod.categorias (categorias_id) ON DELETE CASCADE,
        subcategorias_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        subcategorias_data_atualizacao TIMESTAMP NULL,
        subcategorias_data_inativacao TIMESTAMP NULL,
        subcategorias_status VARCHAR(50) NOT NULL DEFAULT 'ativo'
    );

INSERT INTO
    prod.subcategorias (subcategorias_name, categorias_id)
VALUES
    ('Refrigerantes', 1),
    ('Sucos', 1),
    ('Águas', 1),
    ('Cervejas', 1),
    ('Destilados', 1),
    ('Salgadinhos', 2),
    ('Sanduíches', 2),
    ('Snacks Saudáveis', 2),
    ('Barrinhas', 2),
    ('Chocolates', 3),
    ('Balas e Gomas', 3),
    ('Bolos', 3),
    ('Sorvetes', 3),
    ('Pudins e Mousses', 3),
    ('Pães', 4),
    ('Bolos Caseiros', 4),
    ('Tortas', 4),
    ('Biscoitos', 4),
    ('Arroz', 5),
    ('Feijão', 5),
    ('Massas', 5),
    ('Farinha', 5),
    ('Óleos', 5),
    ('Temperos', 5),
    ('Carnes Bovinas', 6),
    ('Aves', 6),
    ('Peixes e Frutos do Mar', 6),
    ('Ovos', 6),
    ('Proteínas Vegetais', 6),
    ('Frutas Frescas', 7),
    ('Verduras e Legumes', 7),
    ('Orgânicos', 7),
    ('Temperos Naturais', 7),
    ('Ervas e Brotos', 7),
    ('Produtos de Hortifruti', 7),
    ('Desinfetantes', 8),
    ('Detergentes', 8),
    ('Sabões e Sabonetes', 8),
    ('Esponjas e Panos', 8),
    ('Amaciantes', 8),
    ('Desengordurantes', 8),
    ('Limpa Vidros', 8),
    ('Multiuso', 8);

CREATE TABLE
    IF NOT EXISTS prod.unidades (
        unidade_id SERIAL PRIMARY KEY,
        unidade_descricao VARCHAR(100) NOT NULL,
        unidade_abreviacao VARCHAR(10) NOT NULL UNIQUE,
        unidade_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        unidade_data_atualizacao TIMESTAMP NULL
    );

-- Inserindo algumas unidades de medida comuns
INSERT INTO
    prod.unidades (unidade_descricao, unidade_abreviacao)
VALUES
    ('Unidade', 'un'),
    ('Litro', 'L'),
    ('Mililitro', 'mL'),
    ('Quilograma', 'kg'),
    ('Grama', 'g'),
    ('Metro', 'm'),
    ('Centímetro', 'cm'),
    ('Pacote', 'pct'),
    ('Caixa', 'cx');

CREATE TABLE
    IF NOT EXISTS prod.products (
        products_id SERIAL PRIMARY KEY,
        products_name VARCHAR(200) NOT NULL,
        products_price DECIMAL(10, 2) NOT NULL,
        products_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        products_data_atualizacao TIMESTAMP NULL,
        products_data_inativacao TIMESTAMP NULL,
        products_status VARCHAR(50) NOT NULL DEFAULT 'ativo',
        categorias_id INT NOT NULL REFERENCES prod.categorias (categorias_id) ON DELETE CASCADE,
        subcategorias_id INT REFERENCES prod.subcategorias (subcategorias_id) ON DELETE SET NULL,
        unidade_id INT NOT NULL REFERENCES prod.unidades (unidade_id) ON DELETE CASCADE
    );

-- Inserindo os produtos com categorias
INSERT INTO
    prod.products (
        products_name,
        products_price,
        categorias_id,
        subcategorias_id,
        unidade_id
    )
VALUES
    -- Bebidas (Litros ou Mililitros)
    ('Coca-Cola', 4.50, 1, 2, 2), -- Litro (Refrigerantes)
    ('Suco de Laranja', 3.20, 1, 3, 2), -- Mililitro (Sucos)
    ('Água Mineral', 1.00, 1, 2, 2), -- Litro (Águas)
    ('Cerveja Heineken', 6.80, 1, 3, 2), -- Mililitro (Cervejas)
    ('Refrigerante Guaraná', 3.90, 1, 2, 2), -- Litro (Refrigerantes)
    ('Chá Gelado', 2.80, 1, 3, 2), -- Mililitro (Sucos)
    -- Lanches (Unidade)
    ('X-Tudo', 8.90, 2, 1, 1), -- (Salgadinhos)
    ('Pastel de Queijo', 4.50, 2, 1, 1), -- (Salgadinhos)
    ('Coxinha', 3.00, 2, 1, 1), -- (Salgadinhos)
    ('Hambúrguer', 7.99, 2, 1, 1), -- (Sanduíches)
    ('Misto Quente', 5.00, 2, 1, 1), -- (Sanduíches)
    ('Esfiha', 4.25, 2, 1, 1), -- (Salgadinhos)
    -- Doces e Sobremesas (Grama ou Unidade)
    ('Chocolate ao Leite', 7.50, 3, 5, 5), -- Grama (Chocolates)
    ('Sorvete de Morango', 5.25, 3, 3, 2), -- Mililitro (Sorvetes)
    ('Pudim', 4.75, 3, 1, 1), -- (Bolos)
    ('Brigadeiro', 2.50, 3, 5, 5), -- Grama (Balas e Gomas)
    ('Torta de Limão', 6.90, 3, 1, 1), -- (Bolos)
    ('Bolo de Chocolate', 9.99, 3, 1, 1), -- (Bolos)
    -- Produtos de Padaria (Unidade)
    ('Pão Francês', 0.80, 4, 1, 1), -- (Pães)
    ('Croissant', 3.99, 4, 1, 1), -- (Pães)
    ('Bolo de Cenoura', 10.50, 4, 1, 1), -- (Bolos Caseiros)
    ('Rosquinha', 2.75, 4, 1, 1), -- (Biscoitos)
    ('Pão de Queijo', 4.50, 4, 1, 1), -- (Pães)
    ('Baguete', 5.20, 4, 1, 1), -- (Pães)
    -- Alimentos Básicos (Kg ou Unidade)
    ('Arroz 5kg', 20.00, 5, 4, 4), -- Kg (Arroz)
    ('Feijão 1kg', 9.50, 5, 4, 4), -- Kg (Feijão)
    ('Macarrão Espaguete', 3.25, 5, 1, 1), -- (Massas)
    ('Óleo de Soja', 7.99, 5, 2, 2), -- Litro (Óleos)
    ('Sal', 2.30, 5, 4, 4), -- Kg (Temperos)
    ('Açúcar', 4.00, 5, 4, 4), -- Kg (Açúcar)
    -- Carnes e Proteínas (Kg)
    ('Filé de Frango 1kg', 18.50, 6, 4, 4), -- Kg (Carnes Bovinas)
    ('Carne Moída 1kg', 27.90, 6, 4, 4), -- Kg (Carnes Bovinas)
    ('Salmão 500g', 32.80, 6, 5, 5), -- Grama (Peixes e Frutos do Mar)
    ('Linguiça Toscana', 15.99, 6, 4, 4), -- Kg (Carnes Bovinas)
    ('Bife de Contra Filé', 29.50, 6, 4, 4), -- Kg (Carnes Bovinas)
    ('Peito de Peru', 25.75, 6, 4, 4), -- Kg (Carnes Bovinas)
    -- Frutas e Verduras (Kg ou Unidade)
    ('Maçã', 4.20, 7, 4, 4), -- Kg (Frutas Frescas)
    ('Tomate', 5.30, 7, 4, 4), -- Kg (Legumes)
    ('Alface', 2.50, 7, 7, 7), -- Metro (folhas) (Verduras Folhosas)
    ('Banana', 3.00, 7, 4, 4), -- Kg (Frutas Frescas)
    ('Abacaxi', 7.20, 7, 1, 1), -- Unidade (Frutas Frescas)
    ('Cenoura', 4.80, 7, 4, 4), -- Kg (Legumes)
    -- Produtos de Limpeza (Pacote, Unidade ou Litro)
    ('Detergente', 3.60, 8, 2, 2), -- Litro (Detergentes)
    ('Sabão em Pó', 15.90, 8, 4, 4), -- Kg (Limpeza Geral)
    ('Papel Higiênico 12un', 18.00, 8, 8, 8), -- Pacote (Limpeza de Roupas)
    ('Amaciante', 12.50, 8, 2, 2), -- Litro (Limpeza Geral)
    ('Álcool 70%', 8.99, 8, 2, 2), -- Litro (Desinfetantes)
    ('Esponja de Limpeza', 2.75, 8, 1, 1);

--- TABELA DOS STATUS DA LISTA DE COMPRAS
CREATE TABLE
    IF NOT EXISTS prod.lst_compras_status (
        lst_compras_status_id SERIAL PRIMARY KEY,
        lst_compras_status_name VARCHAR(50) NOT NULL,
        lst_compras_status_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

--- DADOS DO STATUS DA LISTA DE COMPRAS
INSERT INTO
    prod.lst_compras_status (lst_compras_status_name)
VALUES
    ('Em Andamento'),
    ('Finalizado'),
    ('Cancelado');

------ TABELA DA LISTA DE COMPRAS
CREATE TABLE
    IF NOT EXISTS prod.lst_compras (
        lst_compras_id SERIAL PRIMARY KEY,
        lst_compras_name VARCHAR(150) NOT NULL,
        lst_compras_valor_total DECIMAL(10, 2) NOT NULL DEFAULT 0,
        lst_compras_total_itens INT NOT NULL DEFAULT 0,
        lst_compras_status_id INT NOT NULL REFERENCES prod.lst_compras_status (lst_compras_status_id) DEFAULT 1,
        lst_compras_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        lst_compras_data_atualizacao TIMESTAMP NULL,
        lst_compras_data_cancelamento TIMESTAMP NULL,
        lst_compras_data_finalizacao TIMESTAMP NULL
    );

-- TABELA DOS ITENS DA LISTA DE COMPRAS
CREATE TABLE
    IF NOT EXISTS prod.lst_compras_itens (
        lst_compras_itens_id SERIAL PRIMARY KEY,
        lst_compras_id INT NOT NULL REFERENCES prod.lst_compras (lst_compras_id) ON DELETE CASCADE,
        products_id INT NOT NULL REFERENCES prod.products (products_id) ON DELETE CASCADE,
        lst_compras_itens_quantidade INT NOT NULL CHECK (lst_compras_itens_quantidade > 0) DEFAULT 0,
        lst_compras_itens_preco DECIMAL(10, 2) NOT NULL CHECK (lst_compras_itens_preco >= 0) DEFAULT 0,
        lst_compras_itens_comprado BOOLEAN NOT NULL DEFAULT FALSE,
        lst_compras_itens_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        lst_compras_itens_data_atualizacao TIMESTAMP NULL
    );

-- INSERINDO LISTA DE COMPRAS
INSERT INTO
    prod.lst_compras (
        lst_compras_name,
        lst_compras_valor_total,
        lst_compras_total_itens,
        lst_compras_status_id,
        lst_compras_data_cancelamento,
        lst_compras_data_finalizacao
    )
VALUES
    ('Compras do Mês', 118.00, 3, 1, NULL, NULL),
    (
        'Almoção da Massa',
        256.99,
        9,
        2,
        NULL,
        CURRENT_TIMESTAMP
    ),
    (
        'Compras do Mês 2',
        98.00,
        6,
        3,
        CURRENT_TIMESTAMP,
        NULL
    );

-- INSERINDO ITENS NA LISTA
INSERT INTO
    prod.lst_compras_itens (
        lst_compras_id,
        products_id,
        lst_compras_itens_quantidade,
        lst_compras_itens_preco,
        lst_compras_itens_comprado
    )
VALUES
    --- lista 1
    (1, 1, 2, 5.50, FALSE),
    (1, 2, 1, 80.00, FALSE),
    (1, 3, 1, 30.00, FALSE),
    --- lista 2
    (2, 4, 2, 1.50, TRUE),
    (2, 2, 1, 60.00, FALSE),
    (2, 5, 1, 20.00, TRUE),
    --- lista 3
    (3, 11, 3, 1.50, FALSE),
    (3, 25, 2, 3.00, FALSE),
    (3, 7, 6, 2.00, FALSE);

DO $$
BEGIN
    -- Reseta as sequências para o valor máximo das respectivas tabelas
    PERFORM setval('prod.categorias_categorias_id_seq', COALESCE((SELECT MAX(categorias_id) FROM prod.categorias), 1), true);
    PERFORM setval('prod.subcategorias_subcategorias_id_seq', COALESCE((SELECT MAX(subcategorias_id) FROM prod.subcategorias), 1), true);
    PERFORM setval('prod.unidades_unidade_id_seq', COALESCE((SELECT MAX(unidade_id) FROM prod.unidades), 1), true);
    PERFORM setval('prod.products_products_id_seq', COALESCE((SELECT MAX(products_id) FROM prod.products), 1), true);
    PERFORM setval('prod.lst_compras_status_lst_compras_status_id_seq', COALESCE((SELECT MAX(lst_compras_status_id) FROM prod.lst_compras_status), 1), true);
    PERFORM setval('prod.lst_compras_lst_compras_id_seq', COALESCE((SELECT MAX(lst_compras_id) FROM prod.lst_compras), 1), true);
    PERFORM setval('prod.lst_compras_itens_lst_compras_itens_id_seq', COALESCE((SELECT MAX(lst_compras_itens_id) FROM prod.lst_compras_itens), 1), true);
END;
$$;