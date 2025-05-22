package helper

import (
	"context"
	"log"
	"crud-api-go/config"

	"github.com/go-redis/redis/v8"
)

var (
	// RedisClient será a instância global do cliente Redis
	RedisClient  *redis.Client	
)

// InitRedis inicializa a conexão com o Redis
func InitRedis() {
	appConfig := config.GetAppConfig()

	opt, err := redis.ParseURL(appConfig.RedisURL)
	if err != nil {
		log.Fatalf("Erro ao analisar REDIS_URL: %v", err)
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Falha ao conectar ao Redis: %v", err)
	}

	log.Println("Conexão com Redis estabelecida com sucesso.")
}
