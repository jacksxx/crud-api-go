package config

import (
	//"net/http"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env. Usando variáveis de ambiente do sistema.")
	}
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnvOrFallback("POSTGRES_HOST", "localhost"),
		Port:     getEnvOrFallback("POSTGRES_PORT", "5432"),
		User:     getEnvOrFallback("POSTGRES_USER", "postgres"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   getEnvOrFallback("POSTGRES_DB", "postgres"),
		SSLMode:  getEnvOrFallback("POSTGRES_SSLMODE", "disable"),
	}
}

// getEnvOrFallback retorna o valor da variável ou um fallback se não definido
func getEnvOrFallback(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
