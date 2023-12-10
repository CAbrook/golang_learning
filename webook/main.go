package main

import (
	"strings"
	"time"

	"github.com/CAbrook/golang_learning/internal/repository"
	"github.com/CAbrook/golang_learning/internal/repository/dao"
	"github.com/CAbrook/golang_learning/internal/service"
	"github.com/CAbrook/golang_learning/internal/web"
	"github.com/CAbrook/golang_learning/internal/web/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// npm run dev
func main() {
	db := initDB()
	server := initWebServer()
	InitUserHandler(db, server)
	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
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
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "xxx")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("this is middleware")
	})

	login := &middlewares.LoginMiddlewareBuilder{}
	//存储数据 直接存cookie
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}

func InitUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}
