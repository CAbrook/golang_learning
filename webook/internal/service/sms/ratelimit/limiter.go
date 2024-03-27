package ratelimit

import (
	"context"
	"errors"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	"github.com/CAbrook/golang_learning/pkg/limter"
)

var errLimited = errors.New("触发限流")

type LimitSMSService struct {
	svc     sms.Service
	limiter limter.Limiter
	key     string
}

func (r *LimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewRateLimitSmsService(svc sms.Service, l limter.Limiter, key string) *LimitSMSService {
	return &LimitSMSService{
		svc:     svc,
		limiter: l,
		key:     key,
	}
}
