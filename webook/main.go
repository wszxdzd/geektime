package main

import (
	repository2 "geektime/week02/webook/internal/repository"
	"geektime/week02/webook/internal/repository/dao"
	"geektime/week02/webook/internal/service"
	"geektime/week02/webook/internal/web"
	"geektime/week02/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	db := InitDb()
	server := InitWebServer()
	InitUser(db, server)

	server.Run(":8080")
}

func InitUser(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository2.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRoutes(server)
}

func InitDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
	login := &middleware.LoginMiddlewareBuilder{}
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}
