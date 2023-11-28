package internal

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type redisMutex struct {
	client   *redis.Client
	key      string
	doneKey  string
	timeout  time.Duration
	ttl      time.Duration
	unlocked bool
}

func (r *redisMutex) Lock() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res := r.client.SetNX(ctx, r.key, true, r.ttl)
	err := res.Err()
	val := res.Val()
	if err != nil {
		return err
	}

	// could not acquire lock
	if !val {
		// wait for lock to be released
		return r.client.BLPop(ctx, r.ttl, r.doneKey).Err()
	}

	// lock acquired
	return nil
}

func (r *redisMutex) Unlock() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if r.unlocked {
		return nil
	}

	tx := r.client.TxPipeline()
	tx.Del(ctx, r.key)
	tx.LPush(ctx, r.doneKey, true)
	_, err := tx.Exec(ctx)

	r.unlocked = true

	return err
}

func DistMutex(client *redis.Client, key string) *redisMutex {
	return &redisMutex{
		client:   client,
		key:      key,
		doneKey:  key + ":done",
		ttl:      5 * time.Second,
		timeout:  1 * time.Second,
		unlocked: false,
	}
}
