package sms

import "context"

// Service 发送短信抽象  屏蔽供应商区别
type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
