/**
 * @Time: 2020/5/1 15:02
 * @Author: solacowa@gmail.com
 * @File: wechat
 * @Software: GoLand
 */

package api

type WechatService interface {
	GetAccessToken()
}

type wechat struct {
}

func NewWechatService() {

}
