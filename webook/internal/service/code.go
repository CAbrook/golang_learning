package service

import (
	"context"
	"github.com/CAbrook/golang_learning/internal/service/sms/tencent"
)

type CodeService struct {
	phoneAndCode map[string]string
	biz          string
}

func (s *CodeService) Send(ctx context.Context, biz, phone string) error {
	// todo add client appid sign name
	tencent.NewService(nil, "test", "test")
	return nil
}

func (s *CodeService) Verify(ctx context.Context, biz, phone string) error {
	return nil
}
