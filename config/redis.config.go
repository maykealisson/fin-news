package config

import (
	"os"

	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func NewRedisClient() *redis.Client {
	logger := log.WithField("config", "redis")
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR não encontrada")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // sem senha por padrão
		DB:       0,  // use default DB
	})

	// Testa a conexão
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.WithError(err).Fatal("Não foi possível conectar ao Redis")
	}

	return rdb
}
