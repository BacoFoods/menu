package internal

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func NewRedis(host, port string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res := rdb.Ping(ctx)

	if res.Err() != nil {
		return nil, res.Err()
	}

	return rdb, nil
}

func MustNewRedis(host, port string) *redis.Client {
	rdb, err := NewRedis(host, port)
	if err != nil {
		panic(err)
	}

	return rdb
}
