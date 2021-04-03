/**
 * @Time: 2021/4/3 下午11:10
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package card

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	FindByUserId(ctx context.Context, userId int64) (res []types.CreditCard, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByUserId(ctx context.Context, userId int64) (res []types.CreditCard, err error) {
	err = s.db.Model(&types.CreditCard{}).
		Preload("Bank").
		Where("user_id = ?", userId).
		Order("bank_id DESC").Find(&res).Error
	return
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
