package service

import (
	"context"
	"fmt"
	"github.com/CAbrook/golang_learning/internal/repository"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	"math/rand"
)

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	//在这是否要直接发送验证码
	if err != nil {
		return err
	}
	const codeTplId = "3277"
	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 屏蔽了因验证次数过多的错误
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
