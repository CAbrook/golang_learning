package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/CAbrook/golang_learning/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct{}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" || path == "/users/login_sms/code/send" || path == "/users/login_sms" {
			return
		}
		//根据约定token在Authorization头部
		// bearer
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			//没登陆没有token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			// 没登陆没有token,Authorization中内容是乱传的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(t *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			//解析出错，说明传的token不对是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			//解析出来可能是非法的或者过期的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime := uc.ExpiresAt
		if expireTime.Before(time.Now()) {
			// 解析出来可能是非法的或者过期的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			//后期告警时候要埋点
			println("输出埋点信息")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//剩余过期时间小于50s刷新（1min过期，每10s刷新一次）
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
