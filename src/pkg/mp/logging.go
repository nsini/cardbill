/**
 * @Time : 3/30/21 3:18 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package mp

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *loggingServer) RecentRepayCount(ctx context.Context, userId int64, recent int) (res int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "RecentRepayCount",
			"userId", userId,
			"recent", recent,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.RecentRepayCount(ctx, userId, recent)
}

func (l *loggingServer) UserInfo(ctx context.Context, userId int64) (res userInfo, err error) {
	panic("implement me")
}

func (l *loggingServer) RecordDetail(ctx context.Context, userId, recordId int64) (res recordDetailResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "RecordDetail",
			"userId", userId,
			"recordId", recordId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.RecordDetail(ctx, userId, recordId)
}

func (l *loggingServer) Statistics(ctx context.Context, userId int64) (res statisticResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Statistics",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Statistics(ctx, userId)
}

func (l *loggingServer) BusinessTypes(ctx context.Context) (res []businessTypesResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "BusinessTypes",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.BusinessTypes(ctx)
}

func (l *loggingServer) RecordAdd(ctx context.Context, userId, cardId int64, amount, rate float64, businessType int64, businessName string, swipeTime *time.Time) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "RecordAdd",
			"userId", userId,
			"cardId", cardId,
			"amount", amount,
			"rate", rate,
			"businessType", businessType,
			"businessName", businessName,
			"swipeTime", swipeTime,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.RecordAdd(ctx, userId, cardId, amount, rate, businessType, businessName, swipeTime)
}

func (l *loggingServer) CreditCards(ctx context.Context, userId int64) (res []cardsResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CreditCards",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreditCards(ctx, userId)
}

func (l *loggingServer) Record(ctx context.Context, userId int64, bankId, cardId int64, start, end *time.Time, page, pageSize int) (res []recordResult, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Record",
			"userId", userId,
			"bankId", bankId,
			"cardId", cardId,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Record(ctx, userId, bankId, cardId, start, end, page, pageSize)
}

func (l *loggingServer) MakeToken(ctx context.Context, appKey string) (token string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "MakeToken",
			"appKey", appKey,
			"token", token,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.MakeToken(ctx, appKey)
}

func (l *loggingServer) Login(ctx context.Context, code, iv, rawData, signature, encryptedData, inviteCode string) (res loginResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Login",
			"code", code,
			"iv", iv,
			"rawData", rawData,
			"signature", signature,
			"encryptedData", encryptedData,
			"inviteCode", inviteCode,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Login(ctx, code, iv, rawData, signature, encryptedData, inviteCode)
}

func (l *loggingServer) RecentRepay(ctx context.Context, userId int64, recent int) (res []recentRepayResult, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "RecentRepay",
			"userId", userId,
			"recent", recent,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.RecentRepay(ctx, userId, recent)
}

func (l *loggingServer) AddCreditCard(ctx context.Context, userId int64, cardName string, bankId int64,
	fixedAmount, maxAmount float64, billingDay, cardHolder int, holderType int, tailNumber int64) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "AddCreditCard",
			"userId", userId,
			"cardName", cardName,
			"bankId", bankId,
			"fixedAmount", fixedAmount,
			"maxAmount", maxAmount,
			"billingDay", billingDay,
			"cardHolder", cardHolder,
			"holderType", holderType,
			"tailNumber", tailNumber,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.AddCreditCard(ctx, userId, cardName, bankId,
		fixedAmount, maxAmount, billingDay, cardHolder, holderType, tailNumber)
}

func (l *loggingServer) BankList(ctx context.Context, bankName string) (res []bankResult, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "BankList",
			"bankName", bankName,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.BankList(ctx, bankName)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "bank", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
