package client

import (
	"context"
	"time"

	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *internal.Config) *RedisClient {
	redisdsn := cfg.RedisDSN
	opt, err := redis.ParseURL(redisdsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse redis DSN")
	}
	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx).Err()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping redis")
	}
	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() {
	r.client.Close()
}

func (r *RedisClient) Health() error {
	return r.client.Ping(context.Background()).Err()
}

func (r *RedisClient) Set(key string, value any, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", utils.NewRedisReturnsNilError("Key not found", err)
	}
	if err != nil {
		return "", utils.NewRedisError("Failed to get key", err)
	}
	return val, nil
}

func (r *RedisClient) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
