package repository

import (
	"context"
	"database/sql"
	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/CAbrook/golang_learning/internal/repository/cache"
	"github.com/CAbrook/golang_learning/internal/repository/dao"
	"log"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	UpdateUserInfo(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	GetProfileById(ctx context.Context, userid int64) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewCacheUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CacheUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Password: u.Password,
		Birthday: u.Birthday,
		About:    u.About,
		Nickname: u.Nickname,
		Phone:    sql.NullString{String: u.Phone, Valid: u.Phone != ""},
	}
}

func (repo *CacheUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Phone:    u.Phone.String,
		About:    u.About,
		Birthday: u.Birthday,
	}
}

func (repo *CacheUserRepository) UpdateUserInfo(ctx context.Context, u domain.User) error {
	return repo.dao.Update(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		About:    u.About,
	})
}

func (repo *CacheUserRepository) GetProfileById(ctx context.Context, userid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, userid)
	if err == nil {
		return du, err
	}
	// 假定err有两种可能：key不存在，redis正常；访问redis error
	u, err := repo.dao.GetProfileById(ctx, userid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	//set error 异步输出日志
	//go func() {
	//	err = repo.cache.Set(ctx, repo.toDomain(u))
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()
	err = repo.cache.Set(ctx, du)
	if err != nil {
		log.Println(err)
	}
	return du, err
}

func (repo *CacheUserRepository) GetProfileByIdV2(ctx context.Context, userid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, userid)
	// 假定err有两种可能：key不存在，redis正常；访问redis error
	switch err {
	case nil:
		return du, err
	case cache.ErrKeyNotExist:
		u, err := repo.dao.GetProfileById(ctx, userid)
		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)
		//set error 异步输出日志
		//go func() {
		//	err = repo.cache.Set(ctx, repo.toDomain(u))
		//	if err != nil {
		//		log.Println(err)
		//	}
		//}()
		err = repo.cache.Set(ctx, du)
		if err != nil {
			log.Println(err)
		}
		return du, err
	default:
		return domain.User{}, err
	}
}

func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
