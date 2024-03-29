package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/CAbrook/golang_learning/internal/integration/startup"
	"github.com/CAbrook/golang_learning/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestUserHandler_SignUp(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServer()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		phone    string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "send success",
			before: func(t *testing.T) {
				// 发送成功啥也不用干
			},
			after: func(t *testing.T) {
				// 验证redis中有数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:123456789"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "123456789",
			wantCode: http.StatusOK,
			wantBody: web.Result{Msg: "send success"},
		},
		{
			name: "phone is null",
			before: func(t *testing.T) {
				// 发送成功啥也不用干
			},
			after: func(t *testing.T) {

			},
			phone:    "",
			wantCode: http.StatusOK,
			wantBody: web.Result{Code: 4, Msg: "please input phone number"},
		},
		{
			name: "code send too many",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:123456789"
				err := rdb.Set(ctx, key, "123456", time.Minute*10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证redis中有数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:123456789"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:    "123456789",
			wantCode: http.StatusOK,
			wantBody: web.Result{Code: 4, Msg: "code send too many"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewReader(
				[]byte(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			// exec
			server.ServeHTTP(recorder, req)
			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}
