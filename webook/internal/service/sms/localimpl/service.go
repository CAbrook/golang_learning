package localimpl

import (
	"context"
	"log"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("verify code is: ", args)
	return nil
}
