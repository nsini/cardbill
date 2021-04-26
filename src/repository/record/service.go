/**
 * @Time : 4/1/21 6:01 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package record

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type Middleware func(Service) Service

type RemainingAmount struct {
	Amount  float64
	Arrival float64
}

type Service interface {
	List(ctx context.Context, userId int64, page, pageSize int, bankId int64, cardIds []int64, start, end *time.Time) (res []types.ExpensesRecord, total int, err error)
	Save(ctx context.Context, record *types.ExpensesRecord) (err error)
	SumAmountCards(ctx context.Context, cardIds []int64, t *time.Time) (ra RemainingAmount, err error)
	FindById(ctx context.Context, userId, id int64) (res types.ExpensesRecord, err error)
	FindByBill(ctx context.Context, userId, cardId int64, beginTime, endTime time.Time) (res []types.ExpensesRecord, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByBill(ctx context.Context, userId, cardId int64, beginTime, endTime time.Time) (res []types.ExpensesRecord, err error) {
	err = s.db.Model(&types.ExpensesRecord{}).
		Preload("Business").
		Where("user_id = ?", userId).
		Where("card_id = ?", cardId).
		Where("created_at >= ?", beginTime.Format("2006-01-02 15:04:05")).
		Where("created_at <= ?", endTime.Format("2006-01-02 15:04:05")).
		Order("created_at DESC").
		Find(&res).Error
	return
}

func (s *service) FindById(ctx context.Context, userId, id int64) (res types.ExpensesRecord, err error) {
	err = s.db.Model(&types.ExpensesRecord{}).
		Preload("CreditCard").
		Preload("CreditCard.Bank").
		Preload("Business").
		Where("id = ? && user_id = ?", id, userId).First(&res).Error
	return
}

func (s *service) SumAmountCards(ctx context.Context, cardIds []int64, t *time.Time) (ra RemainingAmount, err error) {
	query := s.db.Model(&types.ExpensesRecord{}).Select("SUM(amount) AS amount, SUM(arrival) AS arrival")
	if t != nil {
		y, m, _ := t.Date()
		query = query.Where("created_at >= ? AND created_at < ?", fmt.Sprintf("%d-%d-01 00:00:00", y, m), fmt.Sprintf("%d-%d-01 00:00:00", y, m+1))
	}
	err = query.Where("card_id in (?)", cardIds).Scan(&ra).Error
	return
}

func (s *service) Save(ctx context.Context, record *types.ExpensesRecord) (err error) {
	return s.db.Model(record).Save(record).Error
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
	query = query.Count(&total).Offset((page - 1) * pageSize).Limit(pageSize)
	err = query.Find(&res).Error
	return
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
