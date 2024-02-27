package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

var ErrKeyNotExist = redis.Nil

func (c UserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (c UserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}