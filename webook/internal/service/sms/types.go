package sms

import "context"

// Service 发送短信抽象  屏蔽供应商区别
// mockgen -source=./webook/internal/service/sms/types.go - package=smsmocks -destination=./webook/internal/service/sms/mocks/sms.mock.go
type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
