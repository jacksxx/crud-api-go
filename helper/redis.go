package helper

import (
	"context"
	"crud-api-go/config"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	// RedisClient será a instância global do cliente Redis
	RedisClient *redis.Client
	RedisCtx    context.Context
	cancelCtx   context.CancelFunc
)

// InitRedis inicializa a conexão com o Redis
func InitRedis() error {
	if RedisClient != nil {
		return nil // já inicializado
	}
	appConfig := config.GetAppConfig()

	opt, err := redis.ParseURL(appConfig.RedisURL)
	if err != nil {
		log.Fatalf("Erro ao analisar REDIS_URL: %v", err)
	}
	// Cria contexto com cancelamento
	RedisCtx, cancelCtx = context.WithCancel(context.Background())

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("Falha ao conectar ao Redis: %v", err)
	}

	log.Println("Conexão com Redis estabelecida com sucesso.")
	return nil
}

// Finaliza o contexto Redis, ex: em shutdown da aplicação
func ShutdownRedis() {
	if cancelCtx != nil {
		cancelCtx()
	}
}