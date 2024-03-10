package main

import (
	"github.com/CAbrook/golang_learning/internal/repository/cache"
	"github.com/CAbrook/golang_learning/internal/service/sms"
	"github.com/CAbrook/golang_learning/internal/service/sms/localimpl"
	"strings"
	"time"

	"github.com/CAbrook/golang_learning/config"
	"github.com/CAbrook/golang_learning/internal/repository"
	"github.com/CAbrook/golang_learning/internal/repository/dao"
	"github.com/CAbrook/golang_learning/internal/service"
	"github.com/CAbrook/golang_learning/internal/web"
	"github.com/CAbrook/golang_learning/internal/web/middlewares"
	"github.com/gin-contrib/cors"
	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// npm run dev
func main() {
	db := initDB()
	server := initWebServer()
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	codeSvc := initCodeService(redisClient)
	InitUserHandler(db, redisClient, codeSvc, server)
	server.Run(":8080")

	// test nginx
	// server := gin.Default()
	// server.GET("/hello", func(ctx *gin.Context) {
	// 	ctx.String(http.StatusOK, "hellow k8s env")
	// })
	// server.Run(":8080")
}

func initDB() *gorm.DB {
	//"root:root@tcp(localhost:13316)/webook"
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		println(config.Config.DB.DSN)
		panic(err)
	}
	db.Debug()
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"}, //通过Authorization头带token
		//这个是允许前端访问后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "localhost") {
				return true
			}
			return true
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("this is middleware")
	})
	//todo 需要配套使用，此处换成JWT之后Login等接口都需要换成JWT实现
	useJWT(server)
	//useSession(server)
	return server
}

func InitUserHandler(db *gorm.DB, redisClient redis.Cmdable, codeSvc *service.CodeService, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(redisClient)
	ur := repository.NewUserRepository(ud, uc)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us, codeSvc)
	hdl.RegisterRoutes(server)
}

func initCodeService(redisClient redis.Cmdable) *service.CodeService {
	cc := cache.NewCodeCache(redisClient)
	crepo := repository.NewCodeRepository(cc)
	// todo mock smsService
	return service.NewCodeService(crepo, initSms())
}

func initSms() sms.Service {
	return localimpl.NewService()
}

func useJWT(server *gin.Engine) {
	login := middlewares.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

//func useSession(server *gin.Engine) {
//	login := &middlewares.LoginMiddlewareBuilder{}
//	//存储数据 直接存cookie
//	store := cookie.NewStore([]byte("secret"))
//	//基于内存实现
//	// store := redis.NewStore([]byte("6EPTG3HE4W6GX4NLTSGW9LM5EMBGRXZ9"),
//	// 	[]byte("6EPTG3HE4W6GX4NLTSGW9LM5EMBGRXZ0"))
//	//两个key分别时指身份认证和数据加密（二者加上权限控制就是信息安全中三个核心概念）
//	store, err := redis.NewStore(16, "tcp", config.Config.Redis.Addr, "",
//		[]byte("6EPTG3HE4W6GX4NLTSGW9LM5EMBGRXZ9"), []byte("6EPTG3HE4W6GX4NLTSGW9LM5EMBGRXZ0"))
//	if err != nil {
//		panic(err)
//	}
//	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
//}
