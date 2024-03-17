package web

import (
	"bytes"
	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/CAbrook/golang_learning/internal/service"
	svcmocks "github.com/CAbrook/golang_learning/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// export PATH=$PATH:/root/work/go/bin
// mockgen -source=./webook/internal/service/user.go - package=svcmocks -destination=./webook/internal/service/mocks/user.mock.go
func TestUserEmailPattern(t *testing.T) {

}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		// mock
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 构造请求，预期中的输入
		reqBuilder func(t *testing.T) *http.Request
		// 预期中的输出
		wantCode int
		wantBody string
	}{
		{
			name: "signup success",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123123123@qq.com",
					Password: "Ab1$Cdef!",
				}).Return(nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email":"123123123@qq.com",
"password":"Ab1$Cdef!",
"confirmPassword":"Ab1$Cdef!"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "sign up",
		},
		{
			name: "bind fail",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().Signup(gomock.Any(), domain.User{
				//	Email:    "123123123@qq.com",
				//	Password: "password",
				//}).Return(nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email":"123123123@qq.com",
"password":"password"
"confirmPassword":"password"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: "",
		},
		{
			name: "confirmPassword error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				//userSvc.EXPECT().Signup(gomock.Any(), domain.User{
				//	Email:    "123123123@qq.com",
				//	Password: "Ab1$Cdef!",
				//}).Return(nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email":"123123123@qq.com",
"password":"Ab1$Cdef!",
"confirmPassword":"Ab1$Cdef!n"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "ConfirmPassword error",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)
			server := gin.Default()
			hdl.RegisterRoutes(server)
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

//func TestHttp(t *testing.T) {
//	req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("http body")))
//	assert.NoError(t, err)
//	recorder := httptest.NewRecorder()
//
//}

//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	userSvc := svcmocks.NewMockUserService(ctrl)
//	userSvc.EXPECT().Signup(gomock.Any(), domain.User{Id: 1, Email: "123123@qq,com"}).Return(nil)
//
//}
