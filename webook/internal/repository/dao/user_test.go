package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDrive "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		sqlMock func(t *testing.T) *sql.DB
		ctx     context.Context
		user    User
		wantErr error
	}{
		{
			name: "insert success",
			sqlMock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(123, 123)
				// 传入sql的正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			user: User{
				Id: 1110,
				Email: sql.NullString{
					String: `1215343819@qq.com`,
				},
				Password: "password",
				Phone: sql.NullString{
					String: `123456789`,
				},
				Nickname: "nickName",
			},
			wantErr: nil,
		},
		{
			name: "email error",
			sqlMock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 传入sql的正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnError(&mysqlDrive.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			user: User{
				Id: 1110,
				Email: sql.NullString{
					String: `1215343819@qq.com`,
				},
				Password: "password",
				Phone: sql.NullString{
					String: `123456789`,
				},
				Nickname: "nickName",
			},
			wantErr: ErrDuplicateEmail,
		},
		{
			name: "insert error",
			sqlMock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 传入sql的正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("database error"))
				return db
			},
			ctx: context.Background(),
			user: User{
				Id: 1110,
				Email: sql.NullString{
					String: `1215343819@qq.com`,
				},
				Password: "password",
				Phone: sql.NullString{
					String: `123456789`,
				},
				Nickname: "nickName",
			},
			wantErr: errors.New("database error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDb := tc.sqlMock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDb,
				SkipInitializeWithVersion: true, // 不允许初始化的时候查询version
			}),
				&gorm.Config{
					DisableAutomaticPing:   true, // 不允许自动ping
					SkipDefaultTransaction: true, // 跳过自动commit
				})
			assert.NoError(t, err)
			dao := NewUserDao(db)
			err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
