package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanerius/dungeonforge/pkg/config"
)

type RedisClient struct {
	rc *redis.Client
}

func GetRedis(isMocked bool) *RedisClient {

	conf := config.NewIConfig(isMocked)
	user, _ := conf.ReadKeyString("redis_host")
	pass, _ := conf.ReadKeyString("redis_pass")
	port, _ := conf.ReadKeyString("redis_port")

	hostname := fmt.Sprintf("%s:%s", user, port)

	// Initialize the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     hostname, // use default Addr
		Password: pass,     // no password set
		DB:       0,        // use default DB
	})

	// Set a key with a 60-second expiration
	ctx := context.Background()
	err := rdb.Set(ctx, "myKey", "myValue", 60*time.Second).Err()
	if err != nil {
		panic(err)
	}

	return &RedisClient{
		rc: rdb,
	}
}

func (r *RedisClient) Get(ctx context.Context, k string) (string, error) {
	val, err := r.rc.Get(ctx, k).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisClient) Set(ctx context.Context, k, v string, seconds time.Duration) error {
	return r.rc.Set(ctx, "myKey", "myValue", 60*time.Second).Err()
}
