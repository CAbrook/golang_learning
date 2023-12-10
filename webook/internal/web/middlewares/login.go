package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct{}

// builder 模式
func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			//不需要登陆校验
			return
		}
		sess := sessions.Default(ctx)
		if sess.Get("userid") == nil {
			//中断不要往后执行，不执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
