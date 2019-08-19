/**
 * @Time: 2019-08-17 20:34
 * @Author: solacowa@gmail.com
 * @File: bank
 * @Software: GoLand
 */

package types

import "time"

type Bank struct {
	Id        int64     `gorm:"column:id" json:"id"`
	BankName  string    `gorm:"column:bank_name" json:"bank_name"` // 银行名称
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *Bank) TableName() string {
	return "banks"
}
