package repository

import (
	"context"
	"geektime/week02/webook/internal/domain"
	"geektime/week02/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})

}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
func (repo *UserRepository) UpdateById(ctx *gin.Context, user domain.User) error {

	return repo.dao.Update(ctx, dao.User{
		Id:       user.Id,
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,
	})
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repo *UserRepository) Profile(ctx *gin.Context, uid int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	birthday := time.UnixMilli(u.Birthday)
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: birthday,
		AboutMe:  u.AboutMe,
	}, nil
}
