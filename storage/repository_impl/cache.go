package repository_impl

import (
	"errors"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/repository"
	"gopkg.in/redis.v5"
	"time"
)

type redisCache struct {
	client *redis.Client
}

func (r *redisCache) Put(key string, value interface{}, timeout time.Duration) error {
	err := r.client.Set(key, value, timeout).Err()
	return getDatabaseError(err)
}

func (r *redisCache) Get(key string, value interface{}) error {
	res := r.client.Get(key)

	if res.Err() != nil {
		return getDatabaseError(res.Err())
	}

	ok := res.Scan(value)
	return getDatabaseError(ok)
}

func (r *redisCache) Delete(keys ...string) error {
	return r.client.Del(keys...).Err()
}

func NewRedis(c *config.Redis) (repository.Cache, error) {
	r := &redisCache{}

	if c.Type != "tcp" && c.Type != "unix" {
		return r, &Err.Error{
			Code: Err.Einternal,
			Err:  errors.New("redis connection type must be either 'unix' or 'tcp'")}
	}

	r.client = redis.NewClient(&redis.Options{
		Network:  c.Type,
		Addr:     c.Url,
		Password: c.Password,
	})

	ok := r.client.Ping().Val()
	if ok != "PONG" {
		return r, &Err.Error{Code: Err.Einternal, Err: errors.New("unable to connect to redis-server")}
	}

	return r, nil
}
