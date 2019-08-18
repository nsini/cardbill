/**
 * @Time: 2019-08-17 20:33
 * @Author: solacowa@gmail.com
 * @File: creditcard
 * @Software: GoLand
 */

package types

import "time"

type CreditCard struct {
	Id          int64     `gorm:"column:id" json:"id"`
	CardName    string    `gorm:"card_name" json:"card_name"`
	BankId      int64     `gorm:"bank_id" json:"bank_id"`
	FixedAmount float64   `gorm:"fixed_amount" json:"fixed_amount"`
	MaxAmount   float64   `gorm:"max_amount" json:"max_amount"`
	BillingDay  int       `gorm:"billing_day" json:"billing_day"`
	Cardholder  int       `gorm:"cardholder" json:"cardholder"`
	UserId      int64     `gorm:"user_id" json:"user_id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}
