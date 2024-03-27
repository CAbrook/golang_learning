package failover

import (
	"context"
	"errors"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	smsmocks "github.com/CAbrook/golang_learning/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailOverSmsService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mocks   func(ctrl *gomock.Controller) []sms.Service
		wantErr error
	}{
		{
			name: "once success",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			wantErr: nil,
		},
		{
			name: "2th success",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send fail"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			wantErr: nil,
		},
		{
			name: "all fail",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send fail"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send fail"))
				return []sms.Service{svc0, svc1}
			},
			wantErr: errors.New("轮询了所有服务仍未成功"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailOverSmsService(tc.mocks(ctrl))
			err := svc.Send(context.Background(), "test", []string{"123"}, "123")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
