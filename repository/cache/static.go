package cache

import (
	"context"
	"fmt"
	"github.com/MuxiKeStack/be-static/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type StaticCache interface {
	GetStatic(ctx context.Context, name string) (domain.Static, error)
	SetStatic(ctx context.Context, static domain.Static) error
}

type RedisStaticCache struct {
	cmd redis.Cmdable
}

func NewRedisStaticCache(cmd redis.Cmdable) StaticCache {
	return &RedisStaticCache{cmd: cmd}
}

func (cache *RedisStaticCache) GetStatic(ctx context.Context, name string) (domain.Static, error) {
	key := cache.staticKey(name)
	content, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.Static{}, err
	}
	return domain.Static{
		Name:    name,
		Content: content,
	}, nil
}

func (cache *RedisStaticCache) SetStatic(ctx context.Context, static domain.Static) error {
	key := cache.staticKey(static.Name)
	return cache.cmd.Set(ctx, key, static.Content, time.Hour*24*7).Err()
}

func (cache *RedisStaticCache) staticKey(name string) string {
	return fmt.Sprintf("kstack:static:%s", name)
}
