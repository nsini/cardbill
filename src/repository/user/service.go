/**
 * @Time : 3/31/21 11:31 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package user

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	FindByUnionId(ctx context.Context, unionId string) (user types.User, err error)
	Save(ctx context.Context, user *types.User) error
}

type service struct {
	db *gorm.DB
}

func (s *service) Save(ctx context.Context, user *types.User) error {
	return s.db.Model(user).Save(user).Error
}

func (s *service) FindByUnionId(ctx context.Context, unionId string) (user types.User, err error) {
	err = s.db.Model(&types.User{}).
		Where("union_id = ?", unionId).
		Or("open_id = ?", unionId).
		First(&user).Error
	return
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
