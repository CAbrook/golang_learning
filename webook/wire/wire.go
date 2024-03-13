//go:build wireinject

package wire

import (
	"github.com/CAbrook/golang_learning/wire/repository"
	"github.com/CAbrook/golang_learning/wire/repository/dao"
	"github.com/google/wire"
)

// wire 核心是使用了抽象语法树，当执行wire之后会使用抽象语法树分析，分析出入参和出参之后就开始排序，排序完成之后会生成一个init func

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return &repository.UserRepository{}
}
