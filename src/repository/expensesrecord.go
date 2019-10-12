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
	List(userId int64, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error)
	RemainingAmount(cardId int64, billingDay time.Time, cardholder time.Time) (ra *RemainingAmount, err error)
	SumAmountCards(cardIds []int64, t *time.Time) (ra *RemainingAmount, err error)
	SumDays() (sumDays []*SumDay, err error)
}

type RemainingAmount struct {
	Amount  float64
	Arrival float64
}

type SumDay struct {
	Day   string
	Amount float64
}

type expenseRecordRepository struct {
	db *gorm.DB
}

func NewExpenseRecordRepository(db *gorm.DB) ExpenseRecordRepository {
	return &expenseRecordRepository{db}
}

//select DATE_FORMAT(created_at,'%Y%-%m-%d') days,SUM(amount) count from expenses_records group by days order by days desc limit 31;

func (c *expenseRecordRepository) SumDays() (sumDays []*SumDay, err error) {
	query := c.db.Model(&types.ExpensesRecord{})
	err = query.Select("DATE_FORMAT(created_at,'%Y%-%m-%d') day,SUM(amount) amount").
		Group("day").
		Order("day desc").
		Limit(31).Scan(&sumDays).Error
	return
}

func (c *expenseRecordRepository) Create(record *types.ExpensesRecord) (err error) {
	return c.db.Save(record).Error
}

func (c *expenseRecordRepository) List(userId int64, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error) {

	err = c.db.Model(&res).Where("user_id = ?", userId).
		Preload("CreditCard", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Bank")
		}).
		Preload("Business").
		Order("created_at DESC").
		Count(&count).Limit(pageSize).Offset(page * pageSize).Find(&res).Error
	return
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
