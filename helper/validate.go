package helper

import (
	"database/sql"
	"fmt"
)

func ValidateUnit(db *sql.DB, unidadeId int) error {
	var count int
	query := `SELECT COUNT(1) FROM prod.unidades WHERE unidade_id = $1`
	err := db.QueryRow(query, unidadeId).Scan(&count)
	if err != nil {
		return err // Se houver erro ao consultar, retornamos o erro
	}
	if count == 0 {
		return fmt.Errorf("unidade não encontrada") // Se count for 0, categoria não existe
	}
	return nil // Categoria existe
}

// ValidateCategory verifica se a categoria existe no banco de dados.
func ValidateCategory(db *sql.DB, categoriaId int) error {
	var count int
	query := `SELECT COUNT(1) FROM prod.categorias WHERE categorias_id = $1`
	err := db.QueryRow(query, categoriaId).Scan(&count)
	if err != nil {
		return err // Se houver erro ao consultar, retornamos o erro
	}
	if count == 0 {
		return fmt.Errorf("categoria não encontrada") // Se count for 0, categoria não existe
	}
	return nil // Categoria existe
}

// ValidateProduct verifica se o producto existe no banco de dados.
func ValidateProduct(tx *sql.Tx, productId int) error {
	var count int
	query := `SELECT COUNT(1) FROM prod.products WHERE products_id = $1`
	err := tx.QueryRow(query, productId).Scan(&count)
	if err != nil {
		return err // Se houver erro ao consultar, retornamos o erro
	}
	if count == 0 {
		return fmt.Errorf("produto não encontrado") // Se count for 0, produto não existe
	}
	return nil // Produto existe
}