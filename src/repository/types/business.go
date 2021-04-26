/**
 * @Time: 2019-08-18 10:48
 * @Author: solacowa@gmail.com
 * @File: bussiness
 * @Software: GoLand
 */

package types

import "time"

type Business struct {
	Id           int64     `gorm:"column:id;primary_key" json:"id"`
	BusinessName string    `gorm:"column:business_name" json:"business_name"` // 商户类型名称
	Code         int64     `gorm:"column:code" json:"code"`                   // MCC码
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m Business) TableName() string {
	return "businesses"
}
