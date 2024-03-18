package service_test

import (
	"context"
	"errors"
	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/CAbrook/golang_learning/internal/repository"
	repmocks "github.com/CAbrook/golang_learning/internal/repository/mocks"
	"github.com/CAbrook/golang_learning/internal/service"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswork(t *testing.T) {
	password := []byte("Ab1$Cdef!")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123123123!"))
	assert.NotNil(t, err)
}

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "login success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123123123@qq.com").Return(domain.User{
					Email:    "123123123@qq.com",
					Password: "$2a$10$FMMtIjCFX52pLmUbcZaiqO4YmF2.NZw0se7z68TqvbETNDe.bCvNa",
					// Phone:    "123123123",
				}, nil)
				return repo
			},
			email:    "123123123@qq.com",
			password: "Ab1$Cdef!",
			wantUser: domain.User{
				Email:    "123123123@qq.com",
				Password: "$2a$10$FMMtIjCFX52pLmUbcZaiqO4YmF2.NZw0se7z68TqvbETNDe.bCvNa",
				// Phone:    "123123123"
			},
			wantErr: nil,
		},
		{
			name: "user not find",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123123123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123123123@qq.com",
			password: "Ab1$Cdef!",
			wantUser: domain.User{},
			wantErr:  service.ErrInvalidUserOrPassword,
		},
		{
			name: "os error",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123123123@qq.com").Return(domain.User{}, errors.New("db error"))
				return repo
			},
			email:    "123123123@qq.com",
			password: "Ab1$Cdef!",
			wantUser: domain.User{},
			wantErr:  errors.New("db error"),
		},
		{
			name: "login success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123123123@qq.com").Return(domain.User{
					Email:    "123123123@qq.com",
					Password: "$2a$10$FMMtIjCFX52pLmUbcZaiqO4YmF2.NZw0se7z68TqvbETNDe.bCvNa",
					// Phone:    "123123123",
				}, nil)
				return repo
			},
			email:    "123123123@qq.com",
			password: "Ab1$Cdef!11",
			wantUser: domain.User{},
			wantErr:  service.ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := service.NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
