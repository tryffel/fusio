package repository_impl

import (
	"errors"
	"github.com/tryffel/fusio/config"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/storage/repository"
	"gopkg.in/redis.v5"
)

type redisCache struct {
	client *redis.Client
}

func (r *redisCache) Put(key string, value interface{}, timeoutMs int) error {
	panic("implement me")
}

func (r *redisCache) Get(key string) (interface{}, error) {
	panic("implement me")
}

func (r *redisCache) Delete(key string) error {
	panic("implement me")
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
