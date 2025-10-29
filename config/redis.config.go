package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func NewRedisClient() *redis.Client {
	logger := log.WithField("config", "redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
