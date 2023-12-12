package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicate uint16 = 1062
		if me.Number == duplicate {
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {

	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) Update(ctx *gin.Context, u User) error {
	var us User
	dao.db.WithContext(ctx).Where("id=?", u.Id).First(&us)
	us.Nickname = u.Nickname
	us.Birthday = u.Birthday
	us.AboutMe = u.AboutMe
	err := dao.db.Save(&us).Error
	return err

}

func (dao *UserDao) FindById(ctx *gin.Context, uid int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", uid).First(&u).Error
	return u, err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Phone    string
	Password string
	Nickname string
	Birthday int64
	AboutMe  string
	Ctime    int64
	Utime    int64
}
