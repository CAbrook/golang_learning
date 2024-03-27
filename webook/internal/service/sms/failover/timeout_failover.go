package failover

import (
	"context"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSmsService struct {
	svcs []sms.Service
	idx  int32
	cnt  int32
	// 连续超过的阈值,只读
	threshold int32
}

func NewTimeoutFailoverSmsService(svcs []sms.Service, threshold int32) *TimeoutFailoverSmsService {
	return &TimeoutFailoverSmsService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailoverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 重置计数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		// 遇到了错误但是不是超时
		// 可以增加可以不增加，如果强调一定是超时可以不增加
		// 如果是EOF之类的错误，可以考虑直接切换
	}
	return err
}
