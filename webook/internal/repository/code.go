package repository

import (
	"context"
	"github.com/CAbrook/golang_learning/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrCodeSendTooMany

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func (c *CacheCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{cache: c}
}
