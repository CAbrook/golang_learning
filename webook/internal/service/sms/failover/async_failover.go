package failover

import (
	"context"
	"errors"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	"log"
	"sync/atomic"
	"time"
)

// SmsRequest represents a pending SMS request with a retry counter
type SmsRequest struct {
	TplId   string
	Args    []string
	Numbers []string
	Retry   int
}

// AsyncFailOverSmsService manages SMS failover and retry logic
type AsyncFailOverSmsService struct {
	svcs          []sms.Service
	idx           uint64
	queue         chan SmsRequest
	maxRetries    int
	retryInterval time.Duration
	isRateLimited func(sms.Service) bool
	isServiceDown func(sms.Service) bool
}

// NewFailOverSmsService creates a new AsyncFailOverSmsService with specified parameters
func NewAsyncFailOverSmsService(svcs []sms.Service, maxRetries int, retryInterval time.Duration, queueSize int,
	isRateLimited func(sms.Service) bool, isServiceDown func(sms.Service) bool) *AsyncFailOverSmsService {
	svc := &AsyncFailOverSmsService{
		svcs:          svcs,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		queue:         make(chan SmsRequest, queueSize),
		isRateLimited: isRateLimited,
		isServiceDown: isServiceDown,
	}
	go svc.asyncSendLoop()
	return svc
}

// asyncSendLoop handles the asynchronous retry logic
func (f *AsyncFailOverSmsService) asyncSendLoop() {
	for req := range f.queue {
		time.Sleep(f.retryInterval)
		if req.Retry < f.maxRetries {
			req.Retry++
			err := f.Send(context.Background(), req.TplId, req.Args, req.Numbers...)
			if err != nil {
				f.queue <- req
			}
		} else {
			log.Println("Max retries reached for request:", req)
		}
	}
}

// Send attempts to send an SMS using available services, falling back to retry queue if necessary
func (f *AsyncFailOverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		if f.isRateLimited(svc) || f.isServiceDown(svc) {
			continue
		}
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}

	req := SmsRequest{TplId: tplId, Args: args, Numbers: numbers}
	f.queue <- req
	return errors.New("failed to send SMS, saved for retry")
}

// SendV1 attempts to send an SMS using available services starting from a specific index, falling back to retry queue if necessary
func (f *AsyncFailOverSmsService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		if f.isRateLimited(svc) || f.isServiceDown(svc) {
			continue
		}
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}

	req := SmsRequest{TplId: tplId, Args: args, Numbers: numbers}
	f.queue <- req
	return errors.New("failed to send SMS, saved for retry")
}

// isRateLimited 检查服务商是否被限流
func IsRateLimited(svc sms.Service) bool {
	// 假设服务商提供了一个名为 IsRateLimited 的方法来判断是否限流
	if rateLimitedSvc, ok := svc.(interface{ IsRateLimited() bool }); ok {
		return rateLimitedSvc.IsRateLimited()
	}
	return false // 如果服务商没有提供 IsRateLimited 方法，默认返回 false
}

// isServiceDown 检查服务商是否崩溃
func IsServiceDown(svc sms.Service) bool {
	// 假设服务商提供了一个名为 IsServiceDown 的方法来判断是否崩溃
	if downSvc, ok := svc.(interface{ IsServiceDown() bool }); ok {
		return downSvc.IsServiceDown()
	}
	return false // 如果服务商没有提供 IsServiceDown 方法，默认返回 false
}
