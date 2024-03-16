//go:build wireinject

package main

import (
	"github.com/CAbrook/golang_learning/internal/repository"
	"github.com/CAbrook/golang_learning/internal/repository/cache"
	"github.com/CAbrook/golang_learning/internal/repository/dao"
	"github.com/CAbrook/golang_learning/internal/service"
	"github.com/CAbrook/golang_learning/internal/web"
	"github.com/CAbrook/golang_learning/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// export PATH=$PATH:/root/work/go/bin wire env path

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		// dao
		dao.NewUserDao,
		// cache
		cache.NewCodeCache, cache.NewUserCache,
		// repository
		repository.NewCacheUserRepository, repository.NewCodeRepository,
		// service
		ioc.InitSms,
		service.NewUserService,
		service.NewCodeService,
		// handler
		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
