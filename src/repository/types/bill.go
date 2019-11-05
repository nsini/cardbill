/**
 * @Time : 2019-09-18 16:27
 * @Author : solacowa@gmail.com
 * @File : bill
 * @Software: GoLand
 */

package types

import "time"

type Bill struct {
	Id           int64      `gorm:"column:id;primary_key" json:"id"`
	CardId       int64      `gorm:"column:card_id" json:"card_id"`
	Amount       float64    `gorm:"column:amount" json:"amount"`
	RepaymentDay time.Time  `gorm:"column:repayment_day;type:date" json:"repayment_day"`
	IsRepay      bool       `gorm:"column:is_repay" json:"is_repay"`
	RepayTime    *time.Time `gorm:"column:repay_time" json:"repay_time"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	CreditCard   CreditCard `gorm:"ForeignKey:id;AssociationForeignKey:card_id" json:"credit_card"`
}

func (m *Bill) TableName() string {
	return "bills"
}
