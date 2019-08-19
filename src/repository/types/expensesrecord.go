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
	CardId       int64     `gorm:"column:card_id" json:"card_id"`             // 你的银行卡id
	BusinessType int64     `gorm:"column:business_type" json:"business_type"` // 商户类型ID 对应businesses表
	BusinessName string    `gorm:"column:business_name" json:"business_name"` // 商户名称 对应用merchant的名称
	Rate         float64   `gorm:"column:rate" json:"rate"`                   // 费率
	Amount       float64   `gorm:"column:amount" json:"amount"`               // 消费金额
	Arrival      float64   `gorm:"column:arrival" json:"arrival"`             // 实际到账
	UserId       int64     `gorm:"column:user_id" json:"user_id"`             // 用户id
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *ExpensesRecord) TableName() string {
	return "expenses_records"
}
