package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

type AppConfig struct {
	AppEnv               string
	JWTKey               string
	JWTExpirationMinutes int
	JWTRefreshHours      int
	CookieDomain         string
	SameSite             http.SameSite
	RedisURL             string
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

// GetAppConfig inicializa as configurações gerais do aplicativo
func GetAppConfig() AppConfig {
	// Obtenha o ambiente e as configurações de expiração de token
	expMinutes, _ := strconv.Atoi(getEnvOrFallback("JWT_EXPIRATION_MINUTES", "15"))
	refreshHours, _ := strconv.Atoi(getEnvOrFallback("JWT_REFRESH_HOURS", "12"))
	appEnv := getEnvOrFallback("APP_ENV", "development")

	// Configuração de Cookie baseada no ambiente
	cookieDomain := getEnvOrFallback("COOKIE_DOMAIN", "")

	sameSite := http.SameSiteLaxMode
	if appEnv == "production" {
		sameSite = http.SameSiteNoneMode
	}
	// Configuração do Redis
	redisURL := getEnvOrFallback("REDIS_URL", "redis://default@localhost:6379")	

	return AppConfig{
		AppEnv:               appEnv,
		JWTKey:               os.Getenv("JWT_KEY"),
		JWTExpirationMinutes: expMinutes,
		JWTRefreshHours:      refreshHours,
		CookieDomain:         cookieDomain,
		SameSite:             sameSite,
		RedisURL:             redisURL,
	}
}

// getEnvOrFallback retorna o valor da variável ou um fallback se não definido
func getEnvOrFallback(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
