CREATE SCHEMA IF NOT EXISTS prod;

SET
    timezone = 'America/Bahia';

-- -- Função para atualizar a data de modificação (data_atualizacao)
-- CREATE OR REPLACE FUNCTION atualizar_data_atualizacao() 
-- RETURNS TRIGGER AS $$
-- BEGIN
--     -- Atualiza a data de modificação apenas para tabelas que têm a coluna *_data_atualizacao
--     IF TG_TABLE_NAME = 'prod.categorias' THEN
--         NEW.categorias_data_atualizacao := CURRENT_TIMESTAMP;
--     ELSIF TG_TABLE_NAME = 'prod.products' THEN
--         NEW.products_data_atualizacao := CURRENT_TIMESTAMP;
--     -- Adicione mais condições para outras tabelas conforme necessário
--     END IF;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
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
    IF NOT EXISTS prod.products (
        products_id SERIAL PRIMARY KEY,
        products_name VARCHAR(50) NOT NULL,
        products_price DECIMAL(10, 2) NOT NULL,
        products_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        products_data_atualizacao TIMESTAMP NULL,
        products_data_inativacao TIMESTAMP NULL,
        products_status VARCHAR(50) NOT NULL DEFAULT 'ativo',
        categorias_id INT NOT NULL REFERENCES prod.categorias (categorias_id) ON DELETE CASCADE
    );

-- Inserindo os produtos com categorias
INSERT INTO
    prod.products (products_name, products_price, categorias_id)
VALUES
    -- Bebidas (categoria_id = 1)
    ('Coca-Cola', 4.50, 1),
    ('Suco de Laranja', 3.20, 1),
    ('Água Mineral', 1.00, 1),
    ('Cerveja Heineken', 6.80, 1),
    ('Refrigerante Guaraná', 3.90, 1),
    ('Chá Gelado', 2.80, 1),
    -- Lanches (categoria_id = 2)
    ('X-Tudo', 8.90, 2),
    ('Pastel de Queijo', 4.50, 2),
    ('Coxinha', 3.00, 2),
    ('Hambúrguer', 7.99, 2),
    ('Misto Quente', 5.00, 2),
    ('Esfiha', 4.25, 2),
    -- Doces e Sobremesas (categoria_id = 3)
    ('Chocolate ao Leite', 7.50, 3),
    ('Sorvete de Morango', 5.25, 3),
    ('Pudim', 4.75, 3),
    ('Brigadeiro', 2.50, 3),
    ('Torta de Limão', 6.90, 3),
    ('Bolo de Chocolate', 9.99, 3),
    -- Produtos de Padaria (categoria_id = 4)
    ('Pão Francês', 0.80, 4),
    ('Croissant', 3.99, 4),
    ('Bolo de Cenoura', 10.50, 4),
    ('Rosquinha', 2.75, 4),
    ('Pão de Queijo', 4.50, 4),
    ('Baguete', 5.20, 4),
    -- Alimentos Básicos (categoria_id = 5)
    ('Arroz 5kg', 20.00, 5),
    ('Feijão 1kg', 9.50, 5),
    ('Macarrão Espaguete', 3.25, 5),
    ('Óleo de Soja', 7.99, 5),
    ('Sal', 2.30, 5),
    ('Açúcar', 4.00, 5),
    -- Carnes e Proteínas (categoria_id = 6)
    ('Filé de Frango 1kg', 18.50, 6),
    ('Carne Moída 1kg', 27.90, 6),
    ('Salmão 500g', 32.80, 6),
    ('Linguiça Toscana', 15.99, 6),
    ('Bife de Contra Filé', 29.50, 6),
    ('Peito de Peru', 25.75, 6),
    -- Frutas e Verduras (categoria_id = 7)
    ('Maçã', 4.20, 7),
    ('Tomate', 5.30, 7),
    ('Alface', 2.50, 7),
    ('Banana', 3.00, 7),
    ('Abacaxi', 7.20, 7),
    ('Cenoura', 4.80, 7),
    -- Produtos de Limpeza (categoria_id = 8)
    ('Detergente', 3.60, 8),
    ('Sabão em Pó', 15.90, 8),
    ('Papel Higiênico 12un', 18.00, 8),
    ('Amaciante', 12.50, 8),
    ('Álcool 70%', 8.99, 8),
    ('Esponja de Limpeza', 2.75, 8);

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
        lst_compras_status_id INT NOT NULL REFERENCES prod.lst_compras_status (lst_compras_status_id),
        lst_compras_data_cadastro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        lst_compras_data_atualizacao TIMESTAMP NULL
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
        lst_compras_status_id
    )
VALUES
    ('Compras do Mês', 118.00, 3, 1),
    ('Almoção da Massa', 256.99, 9, 2),
    ('Compras do Mês 2', 98.00, 6, 3);

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

-- Trigger para a tabela prod.categorias
CREATE TRIGGER atualizar_categoria BEFORE
UPDATE ON prod.categorias FOR EACH ROW EXECUTE FUNCTION atualizar_data_atualizacao ();

-- Trigger para a tabela prod.products
CREATE TRIGGER atualizar_product BEFORE
UPDATE ON prod.products FOR EACH ROW EXECUTE FUNCTION atualizar_data_atualizacao ();