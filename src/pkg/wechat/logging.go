/**
 * @Time: 2020/10/16 22:08
 * @Author: solacowa@gmail.com
 * @File: logging
 * @Software: GoLand
 */

package wechat

import (
	"context"
	core2 "github.com/chanxuehong/wechat/mch/core"
	"github.com/chanxuehong/wechat/mp/core"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingServer struct {
	logger log.Logger
	Service
	traceId string
}

func NewLoggingServer(logger log.Logger, traceId string, s Service) Service {
	logger = log.With(logger, "logging", "wechat")
	return &loggingServer{
		logger:  level.Info(logger),
		Service: s,
		traceId: traceId,
	}
}

func (s *loggingServer) JsSDK(ctx context.Context, link string) (appId, timestamp, nonceStr, signature string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "JsSDK",
			"link", link,
			"appId", appId,
			"timestamp", timestamp,
			"nonceStr", nonceStr,
			"signature", signature,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.JsSDK(ctx, link)
}

func (s *loggingServer) Callback(ctx context.Context) *core.Server {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Callback",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Callback(ctx)
}

func (s *loggingServer) CallbackMch(ctx context.Context) *core2.Server {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "CallbackMch",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.CallbackMch(ctx)
}

func (s *loggingServer) AuthCodeURL(ctx context.Context) (httpUrl string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AuthCodeURL",
			"httpUrl", httpUrl,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.AuthCodeURL(ctx)
}

func (s *loggingServer) AccessToken(ctx context.Context) (token string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "AccessToken",
			"token", token,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.AccessToken(ctx)
}

func (s *loggingServer) UnifiedOrder(ctx context.Context, tradeNo, clientIp, body, detail, openId string, totalFee int64) (resp UnifiedOrderResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "UnifiedOrder",
			"tradeNo", tradeNo,
			"clientIp", clientIp, "body", body, "detail", detail,
			"openId", openId, "totalFee", totalFee,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.UnifiedOrder(ctx, tradeNo, clientIp, body, detail, openId, totalFee)
}

func (s *loggingServer) OrderQuery(ctx context.Context, transactionId, tradeNo string) (res OrderQueryResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "OrderQuery",
			"transactionId", transactionId,
			"tradeNo", tradeNo,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.OrderQuery(ctx, transactionId, tradeNo)
}

func (s *loggingServer) MPLogin(ctx context.Context, code string) (userInfo UserInfo, sessionKey string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "MPLogin",
			"code", code,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.MPLogin(ctx, code)
}

func (s *loggingServer) GetUserInfo(ctx context.Context, openId string) (userInfo *UserInfo, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "GetUserInfo",
			"openId", openId,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.GetUserInfo(ctx, openId)
}
