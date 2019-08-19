/**
 * @Time : 2019-08-19 09:55
 * @Author : solacowa@gmail.com
 * @File : merchant
 * @Software: GoLand
 */

package types

import "time"

// 商户
type Merchant struct {
	Id           int64     `gorm:"column:id" json:"id"`
	MerchantName string    `gorm:"column:merchant_name;comment('商户名')" json:"merchant_name"`
	BusinessId   int64     `gorm:"column:business_id；comment('商户类型')" json:"business_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *Merchant) TableName() string {
	return "merchants"
}
