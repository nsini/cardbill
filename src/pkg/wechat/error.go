/**
 * @Time : 2020/10/27 4:39 PM
 * @Author : solacowa@gmail.com
 * @File : error
 * @Software: GoLand
 */

package wechat

import (
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"io/ioutil"
	"net/http"
)

type wechatError struct {
	logger log.Logger
}

func (w *wechatError) ServeError(write http.ResponseWriter, r *http.Request, err error) {
	b, _ := ioutil.ReadAll(r.Body)
	_ = level.Error(w.logger).Log("Wecaht", "ServeError", "err", err.Error(), "body", string(b))
}

func newError(logger log.Logger) core.ErrorHandler {
	logger = log.With(logger, "logging", "wechatError")
	return &wechatError{logger: logger}
}
