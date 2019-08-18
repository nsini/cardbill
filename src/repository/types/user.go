/**
 * @Time: 2019-08-18 00:08
 * @Author: solacowa@gmail.com
 * @File: user
 * @Software: GoLand
 */

package types

import "time"

type User struct {
	Id        int64     `gorm:"column:id" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Email     string    `gorm:"column:email" json:"email"`
	GithubId  int64     `gorm:"column:github_id" json:"github_id"`
	OpenId    string    `gorm:"column:open_id" json:"open_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
