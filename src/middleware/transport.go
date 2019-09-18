/**
 * @Time : 2019-08-21 11:16
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"github.com/go-kit/kit/transport/http"
	stdhttp "net/http"
)

func Redirect() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		return ctx
	}
}
