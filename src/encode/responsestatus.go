/**
 * @Time: 2020/3/27 17:34
 * @Author: solacowa@gmail.com
 * @File: responsestatus
 * @Software: GoLand
 */

package encode

import (
	"github.com/pkg/errors"
)

type ResStatus string

var ResponseMessage = map[ResStatus]int{
	Invalid:        400,
	InvalidParams:  400,
	ErrParamsPhone: 401,
	ErrSystem:      500,
	ErrNotfound:    404,

	ErrAccountNotFound:    404,
	ErrAccountLogin:       1002,
	ErrAccountLoginIsNull: 1003,
	ErrAccountNotLogin:    501,
	ErrAccountASD:         1004,
	ErrAccountLocked:      1005,

	ErrAuthNotLogin: 501,

	// User模块错误码
	ErrUserToken:          3400,
	ErrUserInfo:           3401,
	ErrUserBindRel:        3402,
	ErrUserNotfound:       3403,
	ErrUserAuthentication: 3404,
	ErrUserCancel:         3405,
}

const (
	// 公共错误信息
	Invalid        ResStatus = "invalid"
	InvalidParams  ResStatus = "请求参数错误"
	ErrNotfound    ResStatus = "不存在"
	ErrBadRoute    ResStatus = "请求路由错误"
	ErrParamsPhone ResStatus = "手机格式不正确"
	ErrLimiter     ResStatus = "太快了,等我一会儿..."

	// 中间件错误信息
	ErrSystem             ResStatus = "系统错误"
	ErrAccountNotLogin    ResStatus = "用户没登录"
	ErrAuthNotLogin       ResStatus = "请先登录"
	ErrAccountLoginIsNull ResStatus = "用户名和密码不能为空"
	ErrAccountLogin       ResStatus = "用户名或密码错误"
	ErrAccountNotFound    ResStatus = "账号不存在"
	ErrAccountASD         ResStatus = "权限验证失败"
	ErrAccountLocked      ResStatus = "用户已被锁定"

	// User模块错误信息
	ErrUserToken          ResStatus = "Token错误"
	ErrUserInfo           ResStatus = "用户信息获取失败"
	ErrUserBindRel        ResStatus = "用户绑定失败"
	ErrUserNotfound       ResStatus = "分享码不正确"
	ErrUserAuthentication ResStatus = "请先实名认证之后再参加活动"
	ErrUserCancel         ResStatus = "是注销用户,无法参与"
)

func (c ResStatus) String() string {
	return string(c)
}

func (c ResStatus) Error() error {
	return errors.New(string(c))
}

func (c ResStatus) Wrap(err error) error {
	return errors.Wrap(err, string(c))
}
