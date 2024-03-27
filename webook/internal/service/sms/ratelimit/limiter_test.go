package ratelimit

import (
	"context"
	"errors"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	smsmocks "github.com/CAbrook/golang_learning/internal/service/sms/mocks"
	"github.com/CAbrook/golang_learning/pkg/limter"
	limitermocks "github.com/CAbrook/golang_learning/pkg/limter/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (sms.Service, limter.Limiter)
		wantErr error
	}{
		{
			name: "no limited",
			mock: func(ctrl *gomock.Controller) (sms.Service, limter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
			wantErr: nil,
		},
		{
			name: "limited",
			mock: func(ctrl *gomock.Controller) (sms.Service, limter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				//svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
			wantErr: errLimited,
		},
		{
			name: "limited os error",
			mock: func(ctrl *gomock.Controller) (sms.Service, limter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("redis limited error"))
				//svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
			wantErr: errors.New("redis limited error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tc.mock(ctrl)
			svc := NewRateLimitSmsService(smsSvc, l, "tencent")
			err := svc.Send(context.Background(), "test", []string{"111"}, "123")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
