/**
 * @Time : 4/6/21 3:23 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package cardbill

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type Middleware func(Service) Service

type Service interface {
	SumByCards(ctx context.Context, cardIds []int64, t *time.Time, repay Repay) (res BillAmount, err error)
	LastBill(ctx context.Context, cardIds []int64, limit int, t *time.Time) (res []types.Bill, err error)
	CountLastBill(ctx context.Context, cardIds []int64, limit int, t *time.Time) (res int, err error)
	FindById(ctx context.Context, id int64) (res types.Bill, err error)
	Save(ctx context.Context, bill *types.Bill) (err error)
}

type Repay int

const (
	RepayTrue Repay = iota
	RepayFalse
	RepayAll
)

type service struct {
	db *gorm.DB
}

func (s *service) Save(ctx context.Context, bill *types.Bill) (err error) {
	return s.db.Model(bill).Save(bill).Error
}

func (s *service) FindById(ctx context.Context, id int64) (res types.Bill, err error) {
	err = s.db.Model(&types.Bill{}).
		Preload("CreditCard").
		Preload("CreditCard.Bank").
		Where("id = ?", id).First(&res).Error
	return
}

func (s *service) CountLastBill(ctx context.Context, cardIds []int64, limit int, t *time.Time) (res int, err error) {
	query := s.db.Model(&types.Bill{}).Where("card_id in (?)", cardIds)
	if t != nil {
		query = query.Where("repayment_day <= ?", t.Format("2006-01-02")).
			Where("is_repay = ?", false).
			Order("repayment_day asc")
	}
	err = query.Order("id desc").Limit(limit).Count(&res).Error
	return
}

func (s *service) LastBill(ctx context.Context, cardIds []int64, limit int, t *time.Time) (res []types.Bill, err error) {
	query := s.db.Model(&types.Bill{}).Where("card_id in (?)", cardIds)
	if t != nil {
		query = query.Where("repayment_day <= ?", t.Format("2006-01-02")).
			//Where("repayment_day >= ?", time.Now().Format("2006-01-02")).
			Where("is_repay = ?", false).
			Preload("CreditCard", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Bank")
			}).Order("repayment_day asc")
	}
	//Where("is_repay = ?", false).
	err = query.Order("id desc").Limit(limit).Find(&res).Error
	return
}

func (s *service) SumByCards(ctx context.Context, cardIds []int64, t *time.Time, repay Repay) (res BillAmount, err error) {
	query := s.db.Model(&types.Bill{}).Select("SUM(amount) AS amount")
	if t != nil {
		y, m, _ := t.Date()
		query = query.Where("created_at >= ? AND created_at < ?", fmt.Sprintf("%d-%d-01 00:00:00", y, m), fmt.Sprintf("%d-%d-01 00:00:00", y, m+1))
	}
	switch repay {
	case RepayTrue:
		query = query.Where("is_repay = true")
	case RepayFalse:
		query = query.Where("is_repay = false")
	}
	err = query.Where("card_id in (?)", cardIds).Scan(&res).Error
	return
}

type BillAmount struct {
	Amount float64
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
