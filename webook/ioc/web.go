package ioc

import (
	"github.com/CAbrook/golang_learning/internal/web"
	"github.com/CAbrook/golang_learning/internal/web/middlewares"
	"github.com/CAbrook/golang_learning/pkg/ginx/middleware/ratelimit"
	"github.com/CAbrook/golang_learning/pkg/limter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		func(ctx *gin.Context) {
			println("this is my middleware")
		},
		(&middlewares.LoginJWTMiddlewareBuilder{}).CheckLogin(),
		ratelimit.NewBuilder(limter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1)).Build(),
	}
}
