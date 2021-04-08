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

type TotalAmount struct {
	Amount    float64
	MaxAmount float64
}
type Service interface {
	FindByUserId(ctx context.Context, userId int64) (res []types.CreditCard, err error)
	FindById(ctx context.Context, userId, cardId int64) (res types.CreditCard, err error)
	Count(ctx context.Context, userId int64) (total int, err error)
	Sum(ctx context.Context, userId int64, state int) (res TotalAmount, err error)
	FindByBankId(ctx context.Context, bankId int64) (res []types.CreditCard, err error)
	Save(ctx context.Context, card *types.CreditCard) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Save(ctx context.Context, card *types.CreditCard) (err error) {
	return s.db.Model(card).Save(card).Error
}

func (s *service) FindByBankId(ctx context.Context, bankId int64) (res []types.CreditCard, err error) {
	err = s.db.Model(&types.CreditCard{}).
		//Preload("Bank").
		Select("DISTINCT card_name").
		Where("bank_id = ?", bankId).
		Order("card_name").
		Find(&res).Error
	return
}

func (s *service) Sum(ctx context.Context, userId int64, state int) (res TotalAmount, err error) {
	err = s.db.Model(&types.CreditCard{}).Select("SUM(fixed_amount) AS amount, SUM(max_amount) as max_amount").
		Where("user_id = ? AND state = ?", userId, state).Scan(&res).Error
	return
}

func (s *service) Count(ctx context.Context, userId int64) (total int, err error) {
	err = s.db.Model(&types.CreditCard{}).Where("user_id = ? AND state = 0", userId).Count(&total).Error
	return
}

func (s *service) FindById(ctx context.Context, userId, cardId int64) (res types.CreditCard, err error) {
	err = s.db.Model(&types.CreditCard{}).
		Preload("Bank").
		Where("id = ? AND user_id = ?", cardId, userId).First(&res).Error
	return
}

func (s *service) FindByUserId(ctx context.Context, userId int64) (res []types.CreditCard, err error) {
	err = s.db.Model(&types.CreditCard{}).
		Preload("Bank").
		Where("user_id = ?", userId).
		Where("state = ?", 0).
		Order("bank_id DESC").Find(&res).Error
	return
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
