/**
 * @Time: 2019-08-18 00:08
 * @Author: solacowa@gmail.com
 * @File: user
 * @Software: GoLand
 */

package types

import "time"

type User struct {
	Id        int64     `gorm:"column:id;primary_key" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Email     string    `gorm:"column:email" json:"email"`
	AuthId    int64     `gorm:"column:auth_id" json:"auth_id"`
	OpenId    string    `gorm:"column:open_id" json:"open_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *User) TableName() string {
	return "users"
}
