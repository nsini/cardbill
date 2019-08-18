/**
 * @Time: 2019-08-17 20:37
 * @Author: solacowa@gmail.com
 * @File: expensesrecord
 * @Software: GoLand
 */

package types

import "time"

type ExpensesRecord struct {
	Id           int64     `gorm:"column:id" json:"id"`
	CardId       int64     `gorm:"column:card_id" json:"card_id"`
	BusinessType int       `gorm:"column:business_type" json:"business_type"`
	BusinessName string    `gorm:"column:business_name" json:"business_name"`
	RateId       int64     `gorm:"column:rate_id" json:"rate_id"`
	Amount       float64   `gorm:"column:amount" json:"amount"`
	Arrival      float64   `gorm:"column:arrival" json:"arrival"`
	UserId       int64     `gorm:"column:user_id" json:"user_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}
