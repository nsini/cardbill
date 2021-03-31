/**
 * @Time: 2020/10/16 22:10
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package wechat

import (
	"context"
	mchcore "github.com/chanxuehong/wechat/mch/core"
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/encode"
	"net/http"
)

type (
	jsSDKRequest struct {
		Link string `json:"url"`
	}

	jsSDKResponse struct {
		AppId     string `json:"appId"`
		Timestamp string `json:"timestamp"`
		NonceStr  string `json:"nonceStr"`
		Signature string `json:"signature"`
	}

	callbackRequest struct {
		r *http.Request
	}

	callbackResponse struct {
		r      *http.Request
		server *core.Server
	}

	callbackMchResponse struct {
		r      *http.Request
		server *mchcore.Server
	}
)

type Endpoints struct {
	JsSDKEndpoint       endpoint.Endpoint
	CallbackEndpoint    endpoint.Endpoint
	AuthCodeURLEndpoint endpoint.Endpoint
	CallbackMchEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		JsSDKEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(jsSDKRequest)
			appId, timestamp, nonceStr, signature, err := s.JsSDK(ctx, req.Link)
			return encode.Response{
				Data: jsSDKResponse{
					AppId:     appId,
					Timestamp: timestamp,
					NonceStr:  nonceStr,
					Signature: signature,
				},
				Error: err,
			}, err
		},
		CallbackEndpoint:    makeCallbackEndpoint(s),
		CallbackMchEndpoint: makeCallbackMchEndpoint(s),
		AuthCodeURLEndpoint: makeAuthCodeURLEndpoint(s),
	}

	for _, m := range dmw["JsSDK"] {
		eps.JsSDKEndpoint = m(eps.JsSDKEndpoint)
	}
	for _, m := range dmw["Callback"] {
		eps.CallbackEndpoint = m(eps.CallbackEndpoint)
	}
	for _, m := range dmw["AuthCodeURL"] {
		eps.AuthCodeURLEndpoint = m(eps.AuthCodeURLEndpoint)
	}
	for _, m := range dmw["CallbackMch"] {
		eps.CallbackMchEndpoint = m(eps.CallbackMchEndpoint)
	}

	return eps
}

func makeCallbackMchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(callbackRequest)
		server := s.CallbackMch(ctx)
		return encode.Response{
			Data: callbackMchResponse{
				r:      req.r,
				server: server,
			},
			Error: err,
		}, err
	}
}

func makeCallbackEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(callbackRequest)
		server := s.Callback(ctx)
		return encode.Response{
			Data: callbackResponse{
				r:      req.r,
				server: server,
			},
			Error: err,
		}, err
	}
}

func makeAuthCodeURLEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		redirectUri, err := s.AuthCodeURL(ctx)

		return encode.Response{
			Data:  redirectUri,
			Error: err,
		}, err
	}
}
