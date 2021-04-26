/**
 * @Time : 2020/11/17 11:06 AM
 * @Author : solacowa@gmail.com
 * @File : mchhandler
 * @Software: GoLand
 */

package wechat

import (
	mchcore "github.com/chanxuehong/wechat/mch/core"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/repository"
)

type handlerMch struct {
	logger     log.Logger
	traceId    string
	server     *mchcore.Server
	repository repository.Repository
}

func (s *handlerMch) ServeMsg(ctx *mchcore.Context) {
	if ctx.Msg["return_code"] != "SUCCESS" || ctx.Msg["result_code"] != "SUCCESS" {
		_ = level.Warn(s.logger).Log("Msg", ctx.Msg["return_code"])
		return
	}
	//outTradeNo := ctx.Msg["out_trade_no"]
	////s.paymentSvc.PutPaymentStatus(context.Background(), outTradeNo)
	//info, err := s.repository.Payment().FindByTradeNo(outTradeNo)
	//if err != nil {
	//	_ = level.Error(s.logger).Log("repository.Payment", "FindByTradeNo", "err", err.Error())
	//	return
	//}
	//
	//logger := log.With(s.logger, "outTradeNo", outTradeNo, "orderId", info.OrderId)
	//
	//t := time.Now()
	//info.PaymentStatus = types.PaymentStatusPaid
	//info.EndTime = &t
	//if err = s.repository.Payment().Update(&info); err != nil {
	//	_ = level.Error(logger).Log("repository.Payment", "Update", "err", err.Error())
	//}
	//
	//info.Order.Status = types.OrderStatusPayFinish
	//if err = s.repository.Order().Update(&info.Order); err != nil {
	//	_ = level.Error(logger).Log("repository.Order", "Update", "err", err.Error())
	//}
	//
	//order, err := s.repository.Order().Find(info.OrderId)
	//if err != nil {
	//	_ = level.Error(logger).Log("repository.Order", "Find", "err", err.Error())
	//	return
	//}
	//
	//// 入user_course表
	//if err = s.repository.UserCourse().Save(&types.UserCourse{
	//	CourseId: order.CourseId,
	//	UserId:   order.UserId,
	//	OrderId:  order.Id,
	//	VideoNum: order.Course.VideoNum,
	//	Price:    order.Amount,
	//}); err != nil {
	//	_ = level.Error(logger).Log("repository.UserCourse", "Save", "err", err.Error())
	//} else {
	//	if videos, err := s.repository.CourseVideo().FindByCourseId(order.CourseId); err == nil {
	//		for _, v := range videos {
	//			if err = s.repository.UserVideo().Save(&types.UserVideo{
	//				CourseId: v.CourseId,
	//				VideoId:  v.Id,
	//				UserId:   order.UserId,
	//			}); err != nil {
	//				_ = level.Error(logger).Log("repository.UserVideo", "Save", "err", err.Error())
	//			}
	//		}
	//	} else {
	//		_ = level.Error(logger).Log("repository.CourseVideo", "FindByCourseId", "err", err)
	//	}
	//}
	_, _ = ctx.ResponseWriter.Write([]byte("success"))
}

func NewHandlerMch(logger log.Logger, traceId string, repository repository.Repository) mchcore.Handler {
	logger = log.With(logger, "service", "handlerMch")
	return &handlerMch{logger: logger, traceId: traceId, repository: repository}
}
