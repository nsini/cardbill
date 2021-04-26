/**
 * @Time : 2020/10/29 5:40 PM
 * @Author : solacowa@gmail.com
 * @File : handler
 * @Software: GoLand
 */

package wechat

import (
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/chanxuehong/wechat/mp/message/callback/request"
	"github.com/chanxuehong/wechat/mp/message/callback/response"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitcache "github.com/icowan/kit-cache"
	"github.com/nsini/cardbill/src/repository"
	"time"
)

type handler struct {
	logger                    log.Logger
	traceId, appId, appSecret string
	server                    *core.Server
	cacheSvc                  kitcache.Service
	repository                repository.Repository
}

func (s *handler) defaultEventHandler(ctx *core.Context) {
	_ = ctx.NoneResponse()
}

func (s *handler) textMsgHandler(ctx *core.Context) {
	msg := request.GetText(ctx.MixedMsg)

	_ = level.Debug(s.logger).Log("fromUserName", msg.FromUserName, "toUserName", msg.ToUserName, "createTime", msg.CreateTime, "content", msg.Content)

	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	//ctx.RawResponse(resp) // 明文回复
	t := time.Now().Unix()
	_ = ctx.AESResponse(resp, t, ctx.Request.URL.Query().Get("nonce"), nil) // aes密文回复
}

func (s *handler) menuClickEventHandler(ctx *core.Context) {
	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	//ctx.RawResponse(resp) // 明文回复
	t := time.Now().Unix()
	_ = ctx.AESResponse(resp, t, ctx.Request.URL.Query().Get("nonce"), nil) // aes密文回复
}

func (s *handler) defaultMsgHandler(ctx *core.Context) {
	_ = level.Debug(s.logger).Log("default", "msg", "body", string(ctx.MsgPlaintext))
	_ = ctx.NoneResponse()
}

func (s *handler) serveMsg(ctx *core.Context) {
	_ = level.Debug(s.logger).Log("default", "msg", "body", string(ctx.MsgPlaintext))
	_ = ctx.NoneResponse()
}

func (s *handler) eventUnSubscribe(ctx *core.Context) {
	//logger := log.With(s.logger, "method", "eventUnSubscribe")
	//event := request.GetUnsubscribeEvent(ctx.MixedMsg)
	//
	//user, err := s.repository.User().FindSave(event.FromUserName)
	//if err == nil {
	//	user.Subscribed = false
	//	if err = s.repository.User().Update(&user); err != nil {
	//		_ = level.Error(logger).Log("repository.User", "Update", "err", err.Error())
	//	}
	//} else {
	//	_ = level.Error(logger).Log("repository.User", "FindSave", "err", err.Error())
	//}
	_ = ctx.NoneResponse()
}

func (s *handler) eventSubscribe(ctx *core.Context) {
	logger := log.With(s.logger, "method", "eventSubscribe")

	event := request.GetSubscribeEvent(ctx.MixedMsg)
	//var err error
	//if _, err = s.repository.User().FindByWxId(event.FromUserName); err != nil {
	//	if err == gorm.ErrRecordNotFound {
	//		err = s.repository.User().Save(&types.User{
	//			WxId:       event.FromUserName,
	//			Subscribed: true,
	//		})
	//	}
	//}
	//if err != nil {
	//	_ = level.Error(logger).Log("err", err.Error())
	//}

	_ = level.Debug(logger).Log("FromUserName", event.FromUserName, "ToUserName", event.ToUserName)
	_ = ctx.NoneResponse()
}

func NewHandler(logger log.Logger, traceId, oriId, appId, appSecret, token, base64AESKey string, cacheSvc kitcache.Service, repository repository.Repository) *core.Server {
	logger = log.With(logger, "service", "handler")

	s := &handler{logger: logger, traceId: traceId, appId: appId, appSecret: appSecret, cacheSvc: cacheSvc, repository: repository}
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(s.defaultMsgHandler)
	mux.DefaultEventHandleFunc(s.defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, s.textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, s.menuClickEventHandler)
	mux.EventHandleFunc(request.EventTypeSubscribe, s.eventSubscribe)
	mux.EventHandleFunc(request.EventTypeUnsubscribe, s.eventSubscribe)

	return core.NewServer(oriId, appId, token, base64AESKey, mux, newError(logger))
}
