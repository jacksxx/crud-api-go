package helper

import (
	"database/sql"
	"fmt"
)

func ValidateUnit(tx *sql.Tx, unidadeId int) error {
	var count int
	query := `SELECT COUNT(1) FROM prod.unidades WHERE unidade_id = $1`
	err := tx.QueryRow(query, unidadeId).Scan(&count)
	if err != nil {
		return err // Se houver erro ao consultar, retornamos o erro
	}
	if count == 0 {
		return fmt.Errorf("unidade de id: %d  não encontrada", unidadeId) // Se count for 0, categoria não existe
	}
	return nil // Categoria existe
}

// ValidateCategory verifica se a categoria existe no banco de dados.
func ValidateCategory(tx *sql.Tx, categoriaId int) error {
	var count int
	query := `SELECT COUNT(1) FROM prod.categorias WHERE categorias_id = $1`
	err := tx.QueryRow(query, categoriaId).Scan(&count)
	if err != nil {
		return err // Se houver erro ao consultar, retornamos o erro
	}
	if count == 0 {
		return fmt.Errorf("categoria de id: %d não encontrada", categoriaId) // Se count for 0, categoria não existe
	}
	return nil // Categoria existe
}

// ValidateSubCategory verifica se a subcategoria existe no banco de dados.
func ValidateSubCategory(tx *sql.Tx, subcategoriaId int) error {
	var count int
	query := `SELECT COUNT(1) FROM prod.subcategorias WHERE subcategorias_id = $1`
	err := tx.QueryRow(query, subcategoriaId).Scan(&count)
	if err != nil {
		return err // Se houver erro ao consultar, retornamos o erro
	}
	if count == 0 {
		return fmt.Errorf("subcategoria de id: %d não encontrada", subcategoriaId) // Se count for 0, categoria não existe
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
		return fmt.Errorf("produto de id: %d não encontrado", productId) // Se count for 0, produto não existe
	}
	return nil // Produto existe
}
