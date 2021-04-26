/**
 * @Time : 2020/10/16 1:50 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"fmt"
	mchcore "github.com/chanxuehong/wechat/mch/core"
	paycore "github.com/chanxuehong/wechat/mch/core"
	"github.com/chanxuehong/wechat/mch/pay"
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/jssdk"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	kitcache "github.com/icowan/kit-cache"
	"github.com/nsini/cardbill/src/encode"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/util/krand"
	"strconv"
	"time"
)

type UserInfo struct {
	OpenId   string `json:"openid"`   // 用户的唯一标识
	Nickname string `json:"nickname"` // 用户昵称
	Sex      int    `json:"sex"`      // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	City     string `json:"city"`     // 普通用户个人资料填写的城市
	Province string `json:"province"` // 用户个人资料填写的省份
	Country  string `json:"country"`  // 国家, 如中国为CN

	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
	// 用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	HeadImageURL string `json:"headimgurl,omitempty"`

	Privilege []string `json:"privilege,omitempty"` // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionId   string   `json:"unionid,omitempty"`   // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
}

type UnifiedOrderResponse struct {
	// 必选返回
	PrepayId  string `xml:"prepay_id"`  // 微信生成的预支付回话标识，用于后续接口调用中使用，该值有效期为2小时
	TradeType string `xml:"trade_type"` // 调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，详细说明见参数规定

	// 下面字段都是可选返回的(详细见微信支付文档), 为空值表示没有返回, 程序逻辑里需要判断
	DeviceInfo string `xml:"device_info"` // 调用接口提交的终端设备号。
	CodeURL    string `xml:"code_url"`    // trade_type 为 NATIVE 时有返回，可将该参数值生成二维码展示出来进行扫码支付
	MWebURL    string `xml:"mweb_url"`    // trade_type 为 MWEB 时有返回
	AppId      string `json:"app_id"`     // 小程序appId
	ApiKey     string `json:"api_key"`    // 商户Key
}

type OrderQueryResponse struct {
	// 必选返回
	TradeState     string    `xml:"trade_state"`      // 交易状态
	TradeStateDesc string    `xml:"trade_state_desc"` // 对当前查询订单状态的描述和下一步操作的指引
	OpenId         string    `xml:"openid"`           // 用户在商户appid下的唯一标识
	TransactionId  string    `xml:"transaction_id"`   // 微信支付订单号
	OutTradeNo     string    `xml:"out_trade_no"`     // 商户系统的订单号，与请求一致。
	TradeType      string    `xml:"trade_type"`       // 调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，MICROPAY，详细说明见参数规定
	BankType       string    `xml:"bank_type"`        // 银行类型，采用字符串类型的银行标识
	TotalFee       int64     `xml:"total_fee"`        // 订单总金额，单位为分
	CashFee        int64     `xml:"cash_fee"`         // 现金支付金额订单现金支付金额，详见支付金额
	TimeEnd        time.Time `xml:"time_end"`         // 订单支付时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则

	// 下面字段都是可选返回的(详细见微信支付文档), 为空值表示没有返回, 程序逻辑里需要判断
	DeviceInfo         string `xml:"device_info"`          // 微信支付分配的终端设备号
	IsSubscribe        *bool  `xml:"is_subscribe"`         // 用户是否关注公众账号
	SubOpenId          string `xml:"sub_openid"`           // 用户在子商户appid下的唯一标识
	SubIsSubscribe     *bool  `xml:"sub_is_subscribe"`     // 用户是否关注子公众账号
	SettlementTotalFee *int64 `xml:"settlement_total_fee"` // 应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType            string `xml:"fee_type"`             // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFeeType        string `xml:"cash_fee_type"`        // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	Detail             string `xml:"detail"`               // 商品详情
	Attach             string `xml:"attach"`               // 附加数据，原样返回
}

type MpConfig struct {
	MpAppId     string // 小程序
	MpAppSecret string
	MchId       string // 商户
	MchApiKey   string
	AppId       string // 公众号
	AppSecret   string
}

type Service interface {
	JsSDK(ctx context.Context, link string) (appId, timestamp, nonceStr, signature string, err error)

	// 微信公众号回调事件
	Callback(ctx context.Context) *core.Server

	// 商户回调
	CallbackMch(ctx context.Context) *mchcore.Server

	// 获取accessToken
	AccessToken(ctx context.Context) (token string, err error)

	// 获取授权地址
	AuthCodeURL(ctx context.Context) (httpUrl string, err error)

	// 下单 商品信息传过来
	UnifiedOrder(ctx context.Context, tradeNo, clientIp, body, detail, openId string, totalFee int64) (resp UnifiedOrderResponse, err error)

	// 查询订单信息
	OrderQuery(ctx context.Context, transactionId, tradeNo string) (res OrderQueryResponse, err error)

	// 小程序登录 // 可以考虑分成两步 登录和获取用户信息
	MPLogin(ctx context.Context, code string) (userInfo UserInfo, sessionKey string, err error)

	// 获取用户信息
	GetUserInfo(ctx context.Context, openId string) (userInfo *UserInfo, err error)

	GetMpConfig(ctx context.Context) *MpConfig
}

type service struct {
	logger     log.Logger
	traceId    string
	server     *core.Server
	mchServer  *mchcore.Server
	cacheSvc   kitcache.Service
	repository repository.Repository
	appDomain  string

	// 商户
	mchId, mchApiKey string

	// 公众号
	appId, appSecret string

	// 小程序
	mpAppId, mpAppSecret string
}

func (s *service) GetMpConfig(ctx context.Context) *MpConfig {
	return &MpConfig{
		MpAppId:     s.mpAppId,
		MpAppSecret: s.mpAppSecret,
		MchId:       s.mchId,
		MchApiKey:   s.mchApiKey,
		AppId:       s.appId,
		AppSecret:   s.appSecret,
	}
}

func (s *service) OrderQuery(ctx context.Context, transactionId, tradeNo string) (res OrderQueryResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "OrderQuery")
	order, err := pay.OrderQuery2(paycore.NewClient(s.mpAppId, s.mchId, s.mchApiKey, nil), &pay.OrderQueryRequest{
		TransactionId: transactionId,
		OutTradeNo:    tradeNo,
	})
	if err != nil {
		_ = level.Error(logger).Log("pay", "OrderQuery2", "err", err.Error())
		return
	}

	res.OpenId = order.OpenId
	res.TradeState = order.TradeState

	return
}

func (s *service) GetUserInfo(ctx context.Context, openId string) (userInfo *UserInfo, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "MPLogin")
	accessToken, err := s.AccessToken(ctx)

	//accessToken := core.NewDefaultAccessTokenServer(s.mpAppId, s.mpAppSecret, nil)
	//rs, err := accessToken.Token()
	//if err != nil {
	//	return
	//}

	info, err := mpoauth2.GetUserInfo(accessToken, openId, "zh_CN", nil)
	if err != nil {
		_ = level.Error(logger).Log("oauth2", "GetUserInfo", "err", err.Error())
		err = encode.ErrWechatOauthGetUserInfo.Error()
		return
	}

	return &UserInfo{
		OpenId:       info.OpenId,
		Nickname:     info.Nickname,
		Sex:          info.Sex,
		City:         info.City,
		Province:     info.Province,
		Country:      info.Country,
		HeadImageURL: info.HeadImageURL,
		Privilege:    info.Privilege,
		UnionId:      info.UnionId,
	}, err
}

func (s *service) MPLogin(ctx context.Context, code string) (userInfo UserInfo, sessionKey string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "MPLogin")
	ep := mpoauth2.NewEndpoint(s.mpAppId, s.mpAppSecret)
	sess, err := mpoauth2.GetSessionWithClient(ep, code, nil)
	if err != nil {
		_ = level.Error(logger).Log("oauth2", "GetSessionWithClient", "err", err.Error())
		err = encode.ErrWechatOauthSession.Error()
		return
	}
	userInfo.UnionId = sess.UnionId
	userInfo.OpenId = sess.OpenId
	sessionKey = sess.SessionKey
	return

	//accessTokenServer := core.NewDefaultAccessTokenServer(s.mpAppId, s.mpAppSecret, nil)
	//accessToken, e := accessTokenServer.Token()
	//if e != nil {
	//	_ = level.Warn(logger).Log("accessTokenServer", "Token", "err", e.Error())
	//	return
	//}
	//info, e := mpoauth2.GetUserInfo(accessToken, sess.OpenId, "zh_CN", nil)
	//if e != nil {
	//	_ = level.Warn(logger).Log("mpoauth2", "GetUserInfo", "err", e.Error())
	//	return
	//}
	//userInfo.Nickname = info.Nickname
	//userInfo.HeadImageURL = info.HeadImageURL
	//userInfo.Country = info.Country
	//userInfo.City = info.City
	//userInfo.Province = info.Province
	//userInfo.Privilege = info.Privilege
	//userInfo.Sex = info.Sex

	return userInfo, sess.SessionKey, nil

}

func (s *service) UnifiedOrder(ctx context.Context, tradeNo, clientIp, body, detail, openId string, totalFee int64) (resp UnifiedOrderResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UnifiedOrder")
	//clientIp = "101.254.182.6" // TODO: 测试定死
	req := pay.UnifiedOrderRequest{
		Body:           body,
		OutTradeNo:     tradeNo,
		TotalFee:       1, // TODO: 测试定死
		SpbillCreateIP: clientIp,
		NotifyURL:      fmt.Sprintf("%s/wechat/callback/mch", s.appDomain),
		TradeType:      "JSAPI",
		SignType:       paycore.SignType_MD5,
		Detail:         detail,
		Attach:         tradeNo,
		TimeStart:      time.Now(),
		TimeExpire:     time.Now().Add(time.Minute * 30),
		OpenId:         openId,
		DeviceInfo:     "WEB",
	}

	res, err := pay.UnifiedOrder2(paycore.NewClient(s.mpAppId, s.mchId, s.mchApiKey, nil), &req)
	if err != nil {
		_ = level.Error(logger).Log("pay.UnifiedOrder2", "err", err.Error())
		return
	}

	resp.PrepayId = res.PrepayId
	resp.TradeType = res.TradeType
	resp.MWebURL = res.MWebURL
	resp.DeviceInfo = res.DeviceInfo
	resp.CodeURL = res.CodeURL
	resp.AppId = s.mpAppId
	resp.ApiKey = s.mchApiKey

	_ = level.Debug(logger).Log("CodeURL", resp.CodeURL, "DeviceInfo", resp.DeviceInfo, "MWebURL", resp.MWebURL, "PrepayId", resp.PrepayId, "TradeType", resp.TradeType)

	return
}

func (s *service) AuthCodeURL(ctx context.Context) (httpUrl string, err error) {
	logger := log.With(s.logger, "method", "AuthCodeURL", s.traceId, ctx.Value(s.traceId))
	_ = level.Debug(logger).Log("a", "b")
	redirectUri := fmt.Sprintf("%s/web/#/", s.appDomain)
	httpUrl = mpoauth2.AuthCodeURL(s.appId, redirectUri, "snsapi_userinfo", uuid.New().String())
	return
}

func (s *service) AccessToken(ctx context.Context) (token string, err error) {
	logger := log.With(s.logger, "method", "AccessToken", s.traceId, ctx.Value(s.traceId))
	type tokenRes struct {
		AccessToken string `json:"access_token"`
	}
	var res tokenRes
	err = s.cacheSvc.GetCall(ctx, "wechat:accountToken", func(key string) (res interface{}, err error) {
		accessToken := core.NewDefaultAccessTokenServer(s.appId, s.appSecret, nil)
		rs, err := accessToken.Token()
		if err != nil {
			return
		}
		return tokenRes{AccessToken: rs}, nil
	}, time.Hour*24*31, &res)
	if err != nil {
		_ = level.Error(logger).Log("cache", "GetCall", "err", err.Error())
		return
	}

	return res.AccessToken, nil
}

func (s *service) Callback(ctx context.Context) *core.Server {
	return s.server
}

func (s *service) CallbackMch(ctx context.Context) *mchcore.Server {
	return s.mchServer
}

func (s *service) JsSDK(ctx context.Context, link string) (appId, timestamp, nonceStr, signature string, err error) {
	logger := log.With(s.logger, "method", "JsSDK", s.traceId, ctx.Value(s.traceId))
	accessToken := core.NewDefaultAccessTokenServer(s.appId, s.appSecret, nil)
	tk, err := jssdk.NewDefaultTicketServer(core.NewClient(accessToken, nil)).Ticket()
	if err != nil {
		_ = level.Error(logger).Log("jssdk", "NewDefaultTicketServer", "err", err.Error())
		return
	}
	timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr = string(krand.Krand(16, krand.KC_RAND_KIND_ALL))
	signature = jssdk.WXConfigSign(tk, nonceStr, timestamp, link)
	appId = s.appId

	return
}

func New(logger log.Logger, traceId string, wechatConfig map[string]string, appDomain string, cacheSvc kitcache.Service, repository repository.Repository) Service {
	logger = log.With(logger, "service", "wechat")

	mchServer := mchcore.NewServer(wechatConfig["mp.app.id"], wechatConfig["mch.id"], wechatConfig["mch.api.key"], NewHandlerMch(logger, traceId, repository), newError(logger))

	return &service{
		logger:      logger,
		traceId:     traceId,
		appId:       wechatConfig["app.id"],
		appSecret:   wechatConfig["app.secret"],
		mchId:       wechatConfig["mch.id"],
		mchApiKey:   wechatConfig["mch.api.key"],
		mpAppId:     wechatConfig["mp.app.id"],
		mpAppSecret: wechatConfig["mp.app.secret"],
		cacheSvc:    cacheSvc,
		repository:  repository,
		appDomain:   appDomain,
		mchServer:   mchServer,
		server: NewHandler(
			logger,
			traceId,
			wechatConfig["ori.id"],
			wechatConfig["app.id"],
			wechatConfig["app.secret"],
			wechatConfig["token"],
			wechatConfig["encoding.aes.key"],
			cacheSvc,
			repository,
		),
	}
}
