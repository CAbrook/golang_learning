package web

import (
	"net/http"

	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/CAbrook/golang_learning/internal/service"
	regexp "github.com/dlclark/regexp2"

	"github.com/gin-gonic/gin"
)

//type UserHandler struct{}

const (
	emailRegexPattern = `^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	svc              *service.UserService
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.GET("/profile", h.Profile)
	ug.POST("/edit", h.Edit)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:password`
		ConfirmPassword string `json:confirmPassword`
	}
	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	print(req.Email)
	isEmail, err := h.emailRegexExp.MatchString(req.Email)
	println(isEmail)
	if err != nil {
		ctx.String(http.StatusOK, "timeout")
		return
	}

	if !isEmail {
		ctx.String(http.StatusOK, "mail error")
		return
	}

	isPassword, err := h.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "timeout")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "password error")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "ConfirmPassword error")
		return
	}
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		ctx.String(http.StatusOK, "sign up")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "email duplicate")
	default:
		ctx.String(http.StatusOK, "system error")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Profile(ctx *gin.Context) {

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:              svc,
	}
}
