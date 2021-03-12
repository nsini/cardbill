/**
 * @Time: 2019-08-18 00:08
 * @Author: solacowa@gmail.com
 * @File: user
 * @Software: GoLand
 */

package types

import "time"

type LoginType int

const (
	LoginTypeGitHub LoginType = 1
	LoginTypeWechat LoginType = 2
	LoginTypeWxMP   LoginType = 3
)

type User struct {
	Id        int64     `gorm:"column:id;primary_key" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Email     string    `gorm:"column:email" json:"email"`
	AuthId    int64     `gorm:"column:auth_id" json:"auth_id"`
	OpenId    string    `gorm:"column:open_id;index" json:"open_id"`
	UnionId   string    `gorm:"column:union_id;index" json:"union_id"`
	LoginType LoginType `gorm:"column:login_type;default:1;index" json:"login_type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *User) TableName() string {
	return "users"
}
