package failover

import (
	"context"
	"errors"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailOverSmsService struct {
	svcs []sms.Service
	// v1
	// 当前服务商下标
	idx uint64
}

func NewFailOverSmsService(svcs []sms.Service) *FailOverSmsService {
	return &FailOverSmsService{svcs: svcs}
}

func (f *FailOverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 这种轮询方案速度较慢，会导致svcs负载不均衡
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("轮询了所有服务仍未成功")
}

// SendV1 起始下标轮询，并且出错也轮询
func (f *FailOverSmsService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// cpu 高速缓存 可能和内存不一致（多cpu）
	// idx 如果有多个cpu，多个cpu修改了idx。读到的idx不一定是哪个，两个cpu都读的时候不一定会读到内存中的idx，有可能是自己高速缓存中的idx
	//idx := f.idx
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	// 迭代length次
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return err
		}
		log.Println(err)
	}
	return errors.New("轮询了所有服务仍未成功")
}
