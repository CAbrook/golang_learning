package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

// LocalCodeCache 本地缓存实现
type LocalCodeCache struct {
	cache      *lru.Cache
	lock       sync.RWMutex
	expiration time.Duration
}

// NewLocalCodeCache 创建一个新的 LocalCodeCache 实例
func NewLocalCodeCache(c *lru.Cache, expiration time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

// Set 设置验证码到缓存
func (l *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	now := time.Now()

	val, ok := l.cache.Get(key)
	if !ok {
		l.cache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		})
		return nil
	}

	itm, ok := val.(codeItem)
	if !ok {
		return errors.New("系统错误")
	}

	if itm.expire.Sub(now) > time.Minute*9 {
		return ErrCodeSendTooMany
	}

	// 更新验证码
	l.cache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})

	return nil
}

// Verify 验证验证码
func (l *LocalCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	key := l.key(biz, phone)
	val, ok := l.cache.Get(key)
	if !ok {
		return false, ErrKeyNotExist
	}

	itm, ok := val.(codeItem)
	if !ok {
		return false, errors.New("系统错误")
	}

	if itm.cnt <= 0 {
		return false, ErrCodeVerifyTooMany
	}

	// 使用写锁来修改验证码的验证次数
	l.lock.Lock()
	defer l.lock.Unlock()
	itm.cnt--

	return itm.code == inputCode, nil
}

// key 生成缓存的键
func (l *LocalCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

// codeItem 缓存项结构
type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}
