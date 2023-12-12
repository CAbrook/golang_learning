package middlewares

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct{}

// builder 模式
func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			//不需要登陆校验
			return
		}
		sess := sessions.Default(ctx)
		userId := sess.Get("userid")
		if sess.Get("userid") == nil {
			//中断不要往后执行，不执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			println("check not pass")
			return
		}

		//确定需要刷新
		now := time.Now()
		const updateTimeKey = "update_time"
		val := sess.Get(updateTimeKey)
		lastUpdateTime, isTimeType := val.(time.Time)
		if val == nil || (!isTimeType) || now.Sub(lastUpdateTime) > time.Minute {
			sess.Set(updateTimeKey, now)
			//gin中sess存储为覆盖式，如果在此处不设置userid，下一次进来只有updatetime没有userid了
			sess.Set("userid", userId)
			err := sess.Save()
			if err != nil {
				println("Refresh error")
			}
		}

	}
}
