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
