package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Only deletes the key if it matches the owner token.
const releaseScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end
`

// Only extends TTL if the key matches the owner token.
const extendScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("expire", KEYS[1], ARGV[2])
else
	return 0
end
`

// DistributedLock provides Redis-based distributed locking.
// Uses token-based ownership to prevent incorrect release after TTL expiration.
type DistributedLock struct {
	client *redis.Client
	key    string
	token  string
	ttl    time.Duration
}

// Key should be unique per job type (e.g., "scheduler:auto-refresh").
// TTL should be longer than the maximum expected job duration.
func NewDistributedLock(redisURL, key string, ttl time.Duration) (*DistributedLock, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &DistributedLock{
		client: client,
		key:    key,
		ttl:    ttl,
	}, nil
}

// Use this to share a Redis connection pool across multiple components.
func NewDistributedLockWithClient(client *redis.Client, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		client: client,
		key:    key,
		ttl:    ttl,
	}
}

// Returns true if lock was acquired, false if another instance holds it.
// Generates a unique token to ensure only this instance can release the lock.
func (l *DistributedLock) TryAcquire(ctx context.Context) (bool, error) {
	token := uuid.New().String()
	ok, err := l.client.SetNX(ctx, l.key, token, l.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx: %w", err)
	}
	if ok {
		l.token = token
	}
	return ok, nil
}

// Use for long-running jobs. Only extends if this instance still owns the lock.
func (l *DistributedLock) Extend(ctx context.Context) error {
	if l.token == "" {
		return fmt.Errorf("lock not held: no token")
	}

	result, err := l.client.Eval(ctx, extendScript, []string{l.key}, l.token, int(l.ttl.Seconds())).Int()
	if err != nil {
		return fmt.Errorf("redis extend: %w", err)
	}
	if result == 0 {
		return fmt.Errorf("lock not held: token mismatch or expired")
	}
	return nil
}

// Only releases if this instance still owns the lock (token matches).
// Safe to call even if lock expired - will not affect other instance's lock.
func (l *DistributedLock) Release(ctx context.Context) error {
	if l.token == "" {
		return nil
	}

	_, err := l.client.Eval(ctx, releaseScript, []string{l.key}, l.token).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("redis release: %w", err)
	}

	l.token = ""
	return nil
}

func (l *DistributedLock) Close() error {
	return l.client.Close()
}
