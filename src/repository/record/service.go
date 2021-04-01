/**
 * @Time : 4/1/21 6:01 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package record

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type Middleware func(Service) Service

type Service interface {
	List(ctx context.Context, userId int64, page, pageSize int, bankId int64, cardIds []int64, start, end *time.Time) (res []types.ExpensesRecord, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) List(ctx context.Context, userId int64, page, pageSize int, bankId int64, cardIds []int64, start, end *time.Time) (res []types.ExpensesRecord, total int, err error) {
	query := s.db.Model(&types.ExpensesRecord{}).Where("user_id = ?", userId).
		Preload("CreditCard", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Bank")
		}).
		Preload("Business").
		Order("created_at DESC")
	if len(cardIds) > 0 {
		query = query.Where("card_id in (?)", cardIds)
	}
	if start != nil {
		query = query.Where("created_at >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		query = query.Where("created_at <= ?", end.Format("2006-01-02"))
	}
	query = query.Count(&total).Limit(pageSize).Offset(page * pageSize)
	err = query.Find(&res).Error
	return
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
