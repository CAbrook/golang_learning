package repository

import (
	"context"
	"database/sql"
	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/CAbrook/golang_learning/internal/repository/cache"
	cachemocks "github.com/CAbrook/golang_learning/internal/repository/cache/mocks"
	"github.com/CAbrook/golang_learning/internal/repository/dao"
	daomocks "github.com/CAbrook/golang_learning/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_GetProfileById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao)
		ctx  context.Context
		uid  int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "find success and cache not find",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().GetProfileById(gomock.Any(), int64(123)).Return(dao.User{
					Id:       123,
					Email:    sql.NullString{String: "123123123@qq.com"},
					Password: "password",
					Phone:    sql.NullString{String: "123123", Valid: true},
					Nickname: "nickname",
					About:    "name",
					Birthday: "123",
					Ctime:    101,
					Utime:    102,
				}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Nickname: "nickname",
					Email:    "123123123@qq.com",
					Password: "password",
					Phone:    "123123",
					About:    "name",
					Birthday: "123",
					Ctime:    time.UnixMilli(101),
				}).Return(nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Nickname: "nickname",
				Email:    "123123123@qq.com",
				Password: "password",
				Phone:    "123123",
				About:    "name",
				Birthday: "123",
				Ctime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
		{
			name: "cache find success",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:       123,
					Nickname: "nickname",
					Email:    "123123123@qq.com",
					Password: "password",
					Phone:    "123123",
					About:    "name",
					Birthday: "123",
					Ctime:    time.UnixMilli(101)}, nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Nickname: "nickname",
				Email:    "123123123@qq.com",
				Password: "password",
				Phone:    "123123",
				About:    "name",
				Birthday: "123",
				Ctime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
		{
			name: "user not find",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().GetProfileById(gomock.Any(), int64(123)).Return(dao.User{}, dao.ErrRecordNotFound)
				return c, d
			},
			uid:      123,
			ctx:      context.Background(),
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, ud := tc.mock(ctrl)
			svc := NewCacheUserRepository(ud, uc)
			user, err := svc.GetProfileById(tc.ctx, tc.uid)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
