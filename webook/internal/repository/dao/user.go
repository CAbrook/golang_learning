package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("email euplicate")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	GetProfileById(ctx context.Context, userId int64) (User, error)
}

type GormUserDao struct {
	db *gorm.DB
}

type User struct {
	Id int64 `gorm:"primaryKey"`
	// sql.nullString 代表可以为空
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
	Nickname string
	About    string
	Birthday string
	Ctime    int64
	Utime    int64
}

func (dao *GormUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrDuplicateEmail
		}
	}
	return err
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{
		db: db,
	}
}

func (dao *GormUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (dao *GormUserDao) Update(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", u.Id).
		Updates(map[string]interface{}{"Nickname": u.Nickname,
			"Birthday": u.Birthday,
			"About":    u.About,
			"utime":    u.Utime}).Error
}

func (dao *GormUserDao) GetProfileById(ctx context.Context, userId int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", userId).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (dao *GormUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (dao *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}
