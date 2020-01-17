/**
 * @Time : 2019-09-18 16:43
 * @Author : solacowa@gmail.com
 * @File : bill
 * @Software: GoLand
 */

package repository

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

var (
	ErrBillExists = errors.New("账单已存在")
)

type BillRepository interface {
	Create(cardId int64, amount float64, repaymentDay time.Time) error
	Repay(cardId int64, amount float64, cardholder time.Time) error
	SumByCards(cardIds []int64, t *time.Time, repay Repay) (res *BillAmount, err error)
	FindByCardIds(cardId []int64, page, pageSize int) (res []*types.Bill, count int64, err error)
	LastBill(cardIds []int64, limit int, t *time.Time) (res []*types.Bill, err error)
	FindByCardIdAndRepaymentDay(cardId int64, repaymentDay time.Time) (res types.Bill, err error)
}

type Repay int

const (
	RepayTrue Repay = iota
	RepayFalse
	RepayAll
)

type billRepository struct {
	db *gorm.DB
}

type BillAmount struct {
	Amount float64
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

func (c *billRepository) LastBill(cardIds []int64, limit int, t *time.Time) (res []*types.Bill, err error) {
	query := c.db.Model(&types.Bill{}).Where("card_id in (?)", cardIds)
	if t != nil {
		query = query.Where("repayment_day <= ?", t.Format("2006-01-02")).
			Where("repayment_day >= ?", time.Now().Format("2006-01-02")).
			Where("is_repay = ?", false).
			Preload("CreditCard", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Bank")
			}).Order("repayment_day asc")
	}
	//Where("is_repay = ?", false).
	err = query.Order("id desc").Limit(limit).Find(&res).Error
	return
}

func (c *billRepository) FindByCardIdAndRepaymentDay(cardId int64, repaymentDay time.Time) (res types.Bill, err error) {
	err = c.db.Model(&types.Bill{}).Where("card_id = ? AND is_repay = true AND repayment_day > ?", cardId, repaymentDay).
		Order("id desc").
		Limit(1).
		First(&res).Error
	return
}

func (c *billRepository) FindByCardIds(cardId []int64, page, pageSize int) (res []*types.Bill, count int64, err error) {
	err = c.db.Model(&types.Bill{}).Where("card_id in (?)", cardId).
		Order("id desc").
		Preload("CreditCard", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Bank")
		}).
		Count(&count).
		Limit(pageSize).
		Offset(page * pageSize).
		Find(&res).Error

	return
}

func (c *billRepository) SumByCards(cardIds []int64, t *time.Time, repay Repay) (res *BillAmount, err error) {
	var rs BillAmount
	query := c.db.Model(&types.Bill{}).Select("SUM(amount) AS amount")
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
	err = query.Where("card_id in (?)", cardIds).Scan(&rs).Error
	return &rs, err
}

func (c *billRepository) Repay(cardId int64, amount float64, cardholder time.Time) error {
	t := time.Now()
	var bill types.Bill

	if err := c.db.Where("card_id = ? AND repayment_day = ?", cardId, cardholder.Format("2006-01-02")).First(&bill).Error; err != nil {
		return err
	}

	if amount > 0 {
		bill.Amount = amount
	}

	bill.IsRepay = true
	bill.RepayTime = &t

	return c.db.Model(&bill).Updates(bill).Error
}

func (c *billRepository) Create(cardId int64, amount float64, repaymentDay time.Time) error {
	if amount == 0 {
		return nil
	}
	genTime := time.Now()
	bill := types.Bill{
		CardId:       cardId,
		Amount:       amount,
		IsRepay:      false,
		RepaymentDay: repaymentDay,
	}

	if err := c.db.Where("card_id = ? AND created_at >= ? AND created_at <= ?",
		cardId, genTime.Format("2006-01-02"),
		time.Unix(genTime.Unix()+86400, 0).Format("2006-01-02")).First(&bill).Error; err != nil {
		return c.db.Save(&bill).Error
	}

	return ErrBillExists
}
