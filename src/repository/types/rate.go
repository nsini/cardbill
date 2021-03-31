/**
 * @Time: 2019-08-18 11:13
 * @Author: solacowa@gmail.com
 * @File: rate
 * @Software: GoLand
 */

package types

import "time"

type Rate struct {
	Id        int64     `gorm:"column:id;primary_key" json:"id"`
	Score     float64   `gorm:"column:score" json:"score"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m Rate) TableName() string {
	return "rates"
}
