/**
 * @Time: 2019-08-17 20:33
 * @Author: solacowa@gmail.com
 * @File: creditcard
 * @Software: GoLand
 */

package types

import "time"

type CreditCard struct {
	Id                int64     `gorm:"column:id;primary_key" json:"id"`
	CardName          string    `gorm:"card_name" json:"card_name"`                  // 信用卡名
	BankId            int64     `gorm:"bank_id" json:"bank_id"`                      // 银行ID
	FixedAmount       float64   `gorm:"fixed_amount" json:"fixed_amount"`            // 固定额度
	MaxAmount         float64   `gorm:"max_amount" json:"max_amount"`                // 临时额度
	BillingDay        int       `gorm:"billing_day" json:"billing_day"`              // 账单日
	Cardholder        int       `gorm:"cardholder" json:"cardholder"`                // 还款日
	UserId            int64     `gorm:"user_id" json:"user_id"`                      // 用户ID
	State             int       `gorm:"column:state;default(0)" json:"state"`        // 状态: 0正常 1锁定
	CardNumber        int64     `gorm:"column:card_number;null" json:"card_number"`  // 信用卡号
	TailNumber        int64     `gorm:"column:tail_number;null"  json:"tail_number"` // 卡号后四位
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
	User              User      `gorm:"ForeignKey:id;AssociationForeignKey:user_id" json:"user,omitempty"`
	Bank              Bank      `gorm:"ForeignKey:id;AssociationForeignKey:bank_id" json:"bank"`
	BillingAmount     float64   `gorm:"-" json:"billing_amount"`
	NextBillingAmount float64   `gorm:"-" json:"next_billing_amount"`
}

func (m *CreditCard) TableName() string {
	return "credit_cards"
}
