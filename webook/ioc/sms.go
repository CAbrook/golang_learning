package ioc

import (
	"github.com/CAbrook/golang_learning/internal/service/sms"
	"github.com/CAbrook/golang_learning/internal/service/sms/localimpl"
)

func InitSms() sms.Service {
	return localimpl.NewService()
}

func initTencentSmsService() sms.Service {
	return nil
}
