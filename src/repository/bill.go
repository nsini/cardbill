/**
 * @Time : 2019-09-18 16:43
 * @Author : solacowa@gmail.com
 * @File : bill
 * @Software: GoLand
 */

package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

var (
	ErrBillExists = errors.New("账单已存在")
)

type BillRepository interface {
	Create(cardId int64, amount float64) error
	Repay(cardId int64, amount float64) error
}

type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

func (c *billRepository) Repay(cardId int64, amount float64) error {
	t := time.Now()
	var bill types.Bill

	if err := c.db.Where("card_id = ? AND created_at >= ? AND created_at <= ?",
		cardId, t.Format("2006-01-02"),
		time.Unix(t.Unix()+86400, 0).Format("2006-01-02")).First(&bill).Error; err != nil {
		return err
	}

	if amount > 0 {
		bill.Amount = amount
	}

	bill.IsRepay = true
	bill.RepayTime = &t

	return c.db.Model(&bill).Updates(bill).Error
}

func (c *billRepository) Create(cardId int64, amount float64) error {
	genTime := time.Now()
	bill := types.Bill{
		CardId:  cardId,
		Amount:  amount,
		IsRepay: false,
	}

	if err := c.db.Where("card_id = ? AND created_at >= ? AND created_at <= ?",
		cardId, genTime.Format("2006-01-02"),
		time.Unix(genTime.Unix()+86400, 0).Format("2006-01-02")).First(&bill).Error; err != nil {
		return c.db.Save(&bill).Error
	}

	return ErrBillExists
}
