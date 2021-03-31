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

	ErrAuthRegisterSendSms:     2501,
	ErrAuthCheckSmsCode:        2502,
	ErrAuthCheckCaptchaCode:    2503,
	ErrAuthPhoneExists:         2504,
	ErrAuthSave:                2505,
	ErrAuthPasswordNotNull:     2506,
	ErrAuthPhoneOrPassword:     2507,
	ErrAuthNotLogin:            501,
	ErrAuthPublicKey:           2508,
	ErrAuthPrivateKey:          2509,
	ErrAuthRsaPassword:         2510,
	ErrAuthSmsToken:            2511,
	ErrAuthAuthNotFound:        2512,
	ErrAuthPasswordUpdate:      2513,
	ErrAuthCheckCaptchaNotnull: 2514,
	ErrAuthMPLogin:             2515,
	ErrAuthMPLoginCode:         2516,

	ErrWechatOauthSession:     2601,
	ErrWechatOauthGetUserInfo: 2602,

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

	ErrAuthRegisterSendSms     ResStatus = "短信发送失败"
	ErrAuthCheckSmsCode        ResStatus = "短信验证码错误"
	ErrAuthCheckCaptchaCode    ResStatus = "图形验证码错误"
	ErrAuthPhoneExists         ResStatus = "手机号已存在"
	ErrAuthSave                ResStatus = "注册失败"
	ErrAuthPhoneOrPassword     ResStatus = "手机号或密码错误"
	ErrAuthPasswordNotNull     ResStatus = "密码不能为空"
	ErrAuthNotLogin            ResStatus = "请先登录"
	ErrAuthPublicKey           ResStatus = "公钥获取失败"
	ErrAuthPrivateKey          ResStatus = "私钥获取失败"
	ErrAuthRsaPassword         ResStatus = "密码解密失败"
	ErrAuthSmsToken            ResStatus = "Token错误,请重试"
	ErrAuthAuthNotFound        ResStatus = "手机号不存在"
	ErrAuthPasswordUpdate      ResStatus = "密码更新失败"
	ErrAuthCheckCaptchaNotnull ResStatus = "图形验证码不能为空"
	ErrAuthMPLogin             ResStatus = "小程序登录失败"
	ErrAuthMPLoginCode         ResStatus = "Code不能为空"

	ErrWechatOauthSession     ResStatus = "获取授权失败"
	ErrWechatOauthGetUserInfo ResStatus = "获取用户信息失败"
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
