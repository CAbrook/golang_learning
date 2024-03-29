package web

import (
	"github.com/CAbrook/golang_learning/internal/domain"
	"github.com/CAbrook/golang_learning/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//type UserHandler struct{}

const (
	emailRegexPattern = `^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

var JWTKey = []byte("6EPTG3HE4W6GX4NLTSGW9LM5EMBGRXZ7")

type UserHandler struct {
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	svc              service.UserService
	codeSvc          service.CodeService
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("users")
	ug.POST("/signup", h.SignUp)
	//ug.POST("/login", h.Login)
	//todo 此处需要把接口全部换成JWT实现
	ug.POST("/login", h.LoginJWT)
	ug.GET("/profile", h.Profile)
	ug.POST("/edit", h.Edit)
	// 手机验证码登录相关功能
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "os error"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "code error please input again"})
		return
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "os error",
		})
		return
	}
	h.setJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{Msg: "login success"})
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "please input phone number",
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "send success",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "code send too many"})
	default:
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "os error"})
	}
	// todo add log
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
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

	isPassword, err := h.passwordRegexExp.MatchString(req.Password)
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
	case service.ErrDuplicateUser:
		ctx.String(http.StatusOK, "email duplicate")
	default:
		ctx.String(http.StatusOK, "system error")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		h.setJWTToken(ctx, u.Id)
		ctx.String(http.StatusOK, "login success")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "username or password error")
	default:
		ctx.String(http.StatusOK, "system error")
	}
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			//30分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
	}
	//自定义头部
	ctx.Header("x-jwt-token", tokenStr)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userid", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 600,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "system error")
			return
		}
		ctx.String(http.StatusOK, "login success")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "username or password error")
	default:
		ctx.String(http.StatusOK, "system error")
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	// type Profile struct {
	// 	Email    string `json:"email"`
	// 	Phone    string `json:"phone"`
	// 	Nickname string `json:"nickname"`
	// 	Birthday string `json:"birthday"`
	// 	AboutMe  string `json:"aboutMe"`
	// }
	// session 实现
	//sess := sessions.Default(ctx)
	//userId := sess.Get("userid")
	// jwt 实现
	uc := ctx.MustGet("user").(UserClaims)
	userId := uc.Uid
	u, err := h.svc.GetProfileById(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "system error"})
		return
	}
	userProfile := gin.H{
		"Email":    u.Email,
		"Phone":    u.Phone,
		"Nickname": u.Nickname,
		"Birthday": u.Birthday,
		"AboutMe":  u.About,
	}

	// 设置响应头部为JSON格式
	ctx.Header("Content-Type", "application/json")

	// 返回JSON数据
	ctx.JSON(http.StatusOK, userProfile)
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//sess := sessions.Default(ctx)
	//userID := sess.Get("userid")
	uc := ctx.MustGet("user").(UserClaims)
	userID := uc.Uid
	err := h.svc.UpdateUserInfo(ctx, domain.User{
		Id:       userID,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		About:    req.AboutMe,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "Edit error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Edit successful"})
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:              svc,
		codeSvc:          codeSvc,
	}
}
