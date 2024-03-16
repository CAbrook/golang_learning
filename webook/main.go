package main

// npm run dev
func main() {
	//db := initDB()
	//server := initWebServer()
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//codeSvc := initCodeService(redisClient)
	//InitUserHandler(db, redisClient, codeSvc, server)

	server := InitWebServer()
	server.Run(":8080")

	// test nginx
	// server := gin.Default()
	// server.GET("/hello", func(ctx *gin.Context) {
	// 	ctx.String(http.StatusOK, "hellow k8s env")
	// })
	// server.Run(":8080")
}

//func useJWT(server *gin.Engine) {
//	login := middlewares.LoginJWTMiddlewareBuilder{}
//	server.Use(login.CheckLogin())
//}

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
