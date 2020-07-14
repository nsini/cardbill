/**
 * @Time: 2019-08-18 00:32
 * @Author: solacowa@gmail.com
 * @File: expensesrecord
 * @Software: GoLand
 */

package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type ExpenseRecordRepository interface {
	Create(record *types.ExpensesRecord) (err error)
	List(userId int64, page, pageSize int, bankId int64, cardIds []int64, start, end *time.Time) (res []*types.ExpensesRecord, count int64, err error)
	ListByCardId(userId, cardId int64, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error)
	RemainingAmount(cardId int64, billingDay time.Time, cardholder time.Time) (ra *RemainingAmount, err error)
	SumAmountCards(cardIds []int64, t *time.Time) (ra *RemainingAmount, err error)
	SumDays(userId int64) (sumDays []*SumDay, err error)
	SumMonth(userId int64) (sumDays []*SumDay, err error)
}

type RemainingAmount struct {
	Amount  float64
	Arrival float64
}

type SumDay struct {
	Day    string
	Amount float64
}

type expenseRecordRepository struct {
	db *gorm.DB
}

func NewExpenseRecordRepository(db *gorm.DB) ExpenseRecordRepository {
	return &expenseRecordRepository{db}
}

func (c *expenseRecordRepository) SumMonth(userId int64) (sumDays []*SumDay, err error) {
	query := c.db.Model(&types.ExpensesRecord{})
	err = query.Select("DATE_FORMAT(created_at,'%Y%-%m') day,SUM(amount) amount").
		Where("user_id = ?", userId).
		Group("day").
		Order("day desc").
		Limit(13).Scan(&sumDays).Error

	return
}

func (c *expenseRecordRepository) SumDays(userId int64) (sumDays []*SumDay, err error) {
	query := c.db.Model(&types.ExpensesRecord{})
	err = query.Select("DATE_FORMAT(created_at,'%Y%-%m-%d') day,SUM(amount) amount").
		Where("user_id = ?", userId).
		Group("day").
		Order("day desc").
		Limit(31).Scan(&sumDays).Error
	return
}

func (c *expenseRecordRepository) Create(record *types.ExpensesRecord) (err error) {
	return c.db.Save(record).Error
}

func (c *expenseRecordRepository) List(userId int64, page, pageSize int, bankId int64, cardIds []int64, start, end *time.Time) (res []*types.ExpensesRecord, count int64, err error) {
	return c.getList(userId, cardIds, page, pageSize, start, end)
}

func (c *expenseRecordRepository) getList(userId int64, cardId []int64, page, pageSize int, start, end *time.Time) (res []*types.ExpensesRecord, count int64, err error) {
	query := c.db.Model(&res).Where("user_id = ?", userId).
		Preload("CreditCard", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Bank")
		}).
		Preload("Business").
		Order("created_at DESC")
	if len(cardId) > 0 {
		query = query.Where("card_id in (?)", cardId)
	}
	if start != nil {
		query = query.Where("created_at >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		query = query.Where("created_at <= ?", end.Format("2006-01-02"))
	}
	query = query.Count(&count).Limit(pageSize).Offset(page * pageSize)
	err = query.Find(&res).Error

	return
}

func (c *expenseRecordRepository) ListByCardId(userId, cardId int64, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error) {
	var ids []int64
	ids = append(ids, cardId)
	return c.getList(userId, ids, page, pageSize, nil, nil)
}

func (c *expenseRecordRepository) RemainingAmount(cardId int64, billingDay time.Time, endBillingDay time.Time) (ra *RemainingAmount, err error) {
	var rs RemainingAmount
	err = c.db.Raw("SELECT SUM(amount) AS amount FROM expenses_records WHERE card_id = ? AND created_at > ? and created_at <= ?",
		cardId, billingDay.Format("2006-01-02 15:04:05"),
		time.Unix(endBillingDay.Unix()+86400, 0).Format("2006-01-02")).Scan(&rs).Error
	return &rs, err
}

func (c *expenseRecordRepository) SumAmountCards(cardIds []int64, t *time.Time) (ra *RemainingAmount, err error) {
	var rs RemainingAmount
	query := c.db.Model(&types.ExpensesRecord{}).Select("SUM(amount) AS amount, SUM(arrival) AS arrival")
	if t != nil {
		y, m, _ := t.Date()
		query = query.Where("created_at >= ? AND created_at < ?", fmt.Sprintf("%d-%d-01 00:00:00", y, m), fmt.Sprintf("%d-%d-01 00:00:00", y, m+1))
	}
	err = query.Where("card_id in (?)", cardIds).Scan(&rs).Error
	return &rs, err
}
