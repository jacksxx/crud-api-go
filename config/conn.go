package config

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

var DB *sql.DB

func Connect() (*sql.DB, error) {
	dbConfig := GetDBConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DBName, dbConfig.Password, dbConfig.SSLMode,
	)
	// Estabelece a conexão com o banco de dados
	DB, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Verifica se a conexão foi bem-sucedida
	if err := DB.Ping(); err != nil {
		log.Fatalf("Erro ao verificar a conexão com o banco de dados: %v", err)
	}
	log.Println("Conexão com o banco de dados estabelecida com sucesso!")

	// Configura o fuso horário para a sessão do banco de dados
	_, err = DB.Exec("SET TIMEZONE = 'America/Bahia'")
	if err != nil {
		log.Printf("Erro ao configurar o fuso horário para a sessão: %v", err)
		return nil, err
	}

	return DB, nil
}
