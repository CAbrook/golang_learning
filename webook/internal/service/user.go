package service

import (
	"context"
	"errors"
	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/CAbrook/golang_learning/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUser         = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("user is not exit or password error")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	UpdateUserInfo(ctx context.Context, u domain.User) error
	GetProfileById(ctx context.Context, userId int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) UpdateUserInfo(ctx context.Context, u domain.User) error {
	return svc.repo.UpdateUserInfo(ctx, u)
}

func (svc *userService) GetProfileById(ctx context.Context, userId int64) (domain.User, error) {
	return svc.repo.GetProfileById(ctx, userId)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// user not find, need registered
	err = svc.repo.Create(ctx, domain.User{Phone: phone})
	// 两种可能一种为唯一索引冲突（phone），另一种为err != nil
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	// err == nil or user already exists
	// 主从延迟 刚插入user 不一定能找到； 插入插的是主库，查询查的是从库
	// 理论上来说此处要强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}
