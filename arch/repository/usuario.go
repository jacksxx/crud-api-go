package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

type UsuarioRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetUsuarios(filters model.UsuarioFilter) ([]model.Usuario, error)
	GetUsuariosById(id int) (model.Usuario, error)
	CreateUsuario(tx *sql.Tx, usuario model.UsuarioPost, hash string, salt string) (model.UsuarioPost, error)
	UpdateUsuario(tx *sql.Tx, usuario model.UsuarioUpdate) (model.UsuarioUpdate, error)
	InativarUsuario(id int, tx *sql.Tx) error
	AtivarUsuario(id int, tx *sql.Tx) error
	CheckEmail(email string, tx *sql.Tx, usuarioId *int) (bool, error)
	CheckUsuario(usuario string, tx *sql.Tx, usuarioId *int) (bool, error)
	CheckUsuarioAtivo(id int, tx *sql.Tx) (bool, error)
	GetSenhaNSalt(id int, tx *sql.Tx) (string, string, error)
	CountUsuario(filters model.UsuarioFilter) (int, error)
	UpdateSenha(id int, hash string, salt string, tx *sql.Tx) error
}

type usuarioRepository struct {
	connection *sql.DB
}

func NewUsuarioRepository(connection *sql.DB) UsuarioRepository {
	return &usuarioRepository{
		connection: connection,
	}
}

func (r *usuarioRepository) BeginTransaction() (*sql.Tx, error) {
	return r.connection.Begin()
}

func (r *usuarioRepository) GetUsuarios(filters model.UsuarioFilter) ([]model.Usuario, error) {
	query := `
		SELECT id, nome, usuario, email, ativo, data_cadastro, data_atualizacao, data_inativacao
        FROM prod.usuarios
	`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Ativo != nil {
		conditions = append(conditions, fmt.Sprintf("ativo = $%d", argIndex))
		args = append(args, filters.Ativo)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (filters.Page - 1) * filters.Limit
	if offset < 0 {
		offset = 0
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, offset)

	rows, err := r.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the results
	var usuariosList []model.Usuario
	for rows.Next() {
		var usuario model.Usuario
		if err := rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Usuario, &usuario.Email, &usuario.Ativo, &usuario.DataCadastro, &usuario.DataAtualizacao, &usuario.DataInativo); err != nil {
			return nil, err
		}
		usuariosList = append(usuariosList, usuario)
	}

	return usuariosList, nil
}

func (r *usuarioRepository) GetUsuariosById(id int) (model.Usuario, error) {
	query, err := r.connection.Prepare(`
		SELECT id, nome, usuario, email, ativo, data_cadastro, data_atualizacao, data_inativacao
        FROM prod.usuarios
		WHERE id = $1
	`)
	if err != nil {
		fmt.Println("Erro ao preparar consulta:", err)
		return model.Usuario{}, err
	}
	defer query.Close()

	var usuario model.Usuario

	err = query.QueryRow(id).Scan(&usuario.Id, &usuario.Nome, &usuario.Usuario, &usuario.Email, &usuario.Ativo, &usuario.DataCadastro, &usuario.DataAtualizacao, &usuario.DataInativo)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nenhum usuário encontrado com ID:", id)
			return model.Usuario{}, nil
		}
		fmt.Println("Erro na consulta ao banco de dados:", err)
		return model.Usuario{}, err
	}
	fmt.Println("Usuário encontrado:", usuario)

	return usuario, nil
}

func (r *usuarioRepository) CreateUsuario(tx *sql.Tx, usuario model.UsuarioPost, hash string, salt string) (model.UsuarioPost, error) {
	var Id int
	query, err := tx.Prepare(`
	INSERT INTO prod.usuarios (nome, usuario, email, senha, salt) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`)
	if err != nil {
		return model.UsuarioPost{}, fmt.Errorf("error ao preparar consulta: %w", err)
	}
	defer query.Close()

	err = query.QueryRow(usuario.Nome, usuario.Usuario, usuario.Email, hash, salt).Scan(&Id)
	if err != nil {
		return model.UsuarioPost{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	usuario.Id = Id

	return usuario, nil
}

func (r *usuarioRepository) UpdateUsuario(tx *sql.Tx, usuario model.UsuarioUpdate) (model.UsuarioUpdate, error) {
	query, err := tx.Prepare(`
		UPDATE prod.usuarios 
		SET nome = $1, usuario = $2, email = $3, data_atualizacao = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id`)
	if err != nil {
		return model.UsuarioUpdate{}, fmt.Errorf("error ao preparar consulta: %w", err)
	}
	defer query.Close()

	err = query.QueryRow(usuario.Nome, usuario.Usuario, usuario.Email, usuario.Id).Scan(&usuario.Id)
	if err != nil {
		return model.UsuarioUpdate{}, fmt.Errorf("erro ao executar insert: %w", err)
	}

	return usuario, nil

}

func (r *usuarioRepository) InativarUsuario(id int, tx *sql.Tx) error {
	query := `
		UPDATE prod.usuarios 
		SET data_inativacao = CURRENT_TIMESTAMP, ativo = false 
		WHERE id = $1
	`

	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *usuarioRepository) AtivarUsuario(id int, tx *sql.Tx) error {
	query := `
		UPDATE prod.usuarios 
		SET data_inativacao = NULL, ativo = true 
		WHERE id = $1
	`

	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (pr *usuarioRepository) UpdateSenha(id int, hash string, salt string, tx *sql.Tx) error {
	query := `
		UPDATE prod.usuarios 
		SET senha = $1, salt = $2 WHERE id = $3	`

	_, err := tx.Exec(query, hash, salt, id)
	return err
}

func (r *usuarioRepository) CountUsuario(filters model.UsuarioFilter) (int, error) {
	query := `SELECT COUNT(*) FROM prod.usuarios`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Ativo != nil {
		conditions = append(conditions, fmt.Sprintf("ativo = $%d", argIndex))
		args = append(args, filters.Ativo)
		argIndex++
	}

	// Add any additional conditions to the WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.connection.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CHECKEMAIL
// Se for false, não há email - Caso true, há email
func (r *usuarioRepository) CheckEmail(email string, tx *sql.Tx, usuarioId *int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM prod.usuarios WHERE email = $1`
	args := []interface{}{email}
	if usuarioId != nil {
		query += ` AND id <> $2`
		args = append(args, *usuarioId)
	}
	query += `)`

	var exists bool
	err := tx.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao checar email: %v", err)
	}

	return exists, nil
}

// CHECKUSUARIO
// Se for false, não há usuario - Caso true, há usuario
func (r *usuarioRepository) CheckUsuario(usuario string, tx *sql.Tx, usuarioId *int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM prod.usuarios WHERE usuario = $1`
	args := []interface{}{usuario}
	if usuarioId != nil {
		query += ` AND id <> $2`
		args = append(args, *usuarioId)
	}
	query += `)`

	var exists bool
	err := tx.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência do nome de usuário: %v", err)
	}

	return exists, nil
}

// CheckUsuarioAtivo
// Se for false, usuário não está ativo - Caso true, usuário está ativo.
func (r *usuarioRepository) CheckUsuarioAtivo(id int, tx *sql.Tx) (bool, error) {
	query := `SELECT ativo FROM prod.usuarios WHERE id = $1`
	var ativo bool
	err := tx.QueryRow(query, id).Scan(&ativo)
	if err != nil {
		return false, err
	}
	return ativo, nil
}

func (r *usuarioRepository) GetSenhaNSalt(id int, tx *sql.Tx) (string, string, error) {
	query := `SELECT senha, salt FROM prod.usuarios WHERE id = $1`
	var hash string
	var salt string

	// Executa a consulta e armazena o hash e o salt.
	err := tx.QueryRow(query, id).Scan(&hash, &salt)
	return hash, salt, err
}
