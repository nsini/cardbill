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
	Id         int64     `gorm:"column:id;primary_key" json:"id"`
	Username   string    `gorm:"column:username" json:"username"`
	Email      string    `gorm:"column:email" json:"email"`
	AuthId     int64     `gorm:"column:auth_id" json:"auth_id"`
	OpenId     string    `gorm:"column:open_id;index" json:"open_id"`
	UnionId    string    `gorm:"column:union_id;index" json:"union_id"`
	LoginType  LoginType `gorm:"column:login_type;default:1;index" json:"login_type"`
	Nickname   string    `gorm:"column:nickname;null" json:"nickname"`      // 用户昵称
	Sex        int       `gorm:"column:sex;null" json:"sex"`                // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	City       string    `gorm:"column:city;null" json:"city"`              // 普通用户个人资料填写的城市
	Province   string    `gorm:"column:province;null" json:"province"`      // 用户个人资料填写的省份
	Country    string    `gorm:"column:country;null" json:"country"`        // 国家, 如中国为CN
	Avatar     string    `gorm:"column:avatar;null;size:500" json:"avatar"` // 头像地址
	Desc       string    `gorm:"column:desc;type:text;null" json:"desc"`    // 描述
	Remark     string    `gorm:"column:remark;null" json:"remark"`          // 标注
	Subscribed bool      `gorm:"column:subscribed;null" json:"subscribed"`  // 是否订阅
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m User) TableName() string {
	return "users"
}
