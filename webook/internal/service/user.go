package service

import (
	"context"
	"errors"
	"geektime/week02/webook/internal/domain"
	"geektime/week02/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvaildUserOrPassword = errors.New("账号或者邮箱错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, user domain.User) error {
	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	//
	return svc.repo.Create(ctx, user)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvaildUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvaildUserOrPassword
	}
	return u, nil

}

func (svc *UserService) Edit(ctx *gin.Context, user domain.User) error {
	return svc.repo.UpdateById(ctx, user)

}

func (svc *UserService) Profile(ctx *gin.Context, uid int64) (domain.User, error) {
	u, err := svc.repo.Profile(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
