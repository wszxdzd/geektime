package web

import (
	"geektime/week02/webook/internal/domain"
	"geektime/week02/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	svc              *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:              svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	isPassword, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	if !isPassword {
		ctx.String(http.StatusOK, "password invalid")
		return
	}

	err = u.svc.Signup(ctx, domain.User{Email: req.Email, Password: req.Password})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}

}

func (u *UserHandler) login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	userinfo, err := u.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", userinfo.Id)
		sess.Options(sessions.Options{
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "服务器异常")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvaildUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码错误")
	default:
		ctx.String(http.StatusOK, "系统错误")

	}
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type UpdateReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req UpdateReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Nickname == "" {
		ctx.String(http.StatusOK, "昵称不能为空")
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {

		ctx.String(http.StatusOK, "日期格式不对")
		return
	}
	// get uid
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	err = u.svc.Edit(ctx, domain.User{
		Id:       userId.(int64),
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "更改错误")
		return
	}
	ctx.String(http.StatusOK, "111")

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Nickname string
		Email    string
		Phone    string
		Birthday time.Time
		AboutMe  string
	}
	sess := sessions.Default(ctx)
	uid := sess.Get("userId")

	du, err := u.svc.Profile(ctx, uid.(int64))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	profile := Profile{
		Nickname: du.Nickname,
		Email:    du.Email,
		Phone:    du.Phone,
		Birthday: du.Birthday,
		AboutMe:  du.AboutMe,
	}
	ctx.JSON(http.StatusOK, profile)

}
