/**
 * @Time: 2019-08-18 09:48
 * @Author: solacowa@gmail.com
 * @File: main
 * @Software: GoLand
 */

package service

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/config"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/encode"
	"github.com/nsini/cardbill/src/logging"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/pkg/auth"
	"github.com/nsini/cardbill/src/pkg/bank"
	"github.com/nsini/cardbill/src/pkg/bill"
	"github.com/nsini/cardbill/src/pkg/business"
	"github.com/nsini/cardbill/src/pkg/creditcard"
	"github.com/nsini/cardbill/src/pkg/dashboard"
	"github.com/nsini/cardbill/src/pkg/merchant"
	"github.com/nsini/cardbill/src/pkg/mp"
	"github.com/nsini/cardbill/src/pkg/record"
	"github.com/nsini/cardbill/src/pkg/user"
	"github.com/nsini/cardbill/src/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/robfig/cron"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger log.Logger

var (
	fs         = flag.NewFlagSet("cardbill", flag.ExitOnError)
	httpAddr   = fs.String("http-addr", ":8080", "HTTP listen address")
	configFile = fs.String("config-file", "app.cfg", "server config file")

	db                 *gorm.DB
	cf                 *config.Config
	tracer             opentracing.Tracer
	appName, namespace string
)

func Run() {
	logger = log.NewLogfmtLogger(log.StdlibWriter{})
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = level.NewFilter(logger, level.AllowAll())

	err := fs.Parse(os.Args[1:])
	if err != nil {
		_ = level.Error(logger).Log("fs", "Parse", "err", err.Error())
		return
	}

	cf, err = config.NewConfig(*configFile)
	if err != nil {
		_ = level.Error(logger).Log("config", "NewConfig", "err", err.Error())
		return
	}

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		cf.GetString(config.SectionMysql, "user"),
		cf.GetString(config.SectionMysql, "password"),
		cf.GetString(config.SectionMysql, "host"),
		cf.GetString(config.SectionMysql, "port"),
		cf.GetString(config.SectionMysql, "database"))

	// 连接数据库
	db, err = mysqlclient.NewMysql(dbUrl, cf.GetBool(config.SectionServer, "debug"))
	if err != nil {
		_ = level.Error(logger).Log("db", "connect", "err", err)
		return
	}

	store := repository.NewRepository(db, logger, logging.TraceId)

	var (
		recordSvc     = record.NewService(logger, store)
		bankSvc       = bank.NewService(logger, store)
		creditCardSvc = creditcard.NewService(logger, store)
		userSvc       = user.NewService(logger, store)
		businessSvc   = business.NewService(logger, store)
		authSvc       = auth.NewService(logger, cf, store)
		billSvc       = bill.NewService(logger, store)
		dashboardSvc  = dashboard.NewService(logger, store)
		merchantSvc   = merchant.NewService(logger, store)
		mpSvc         = mp.New(logger, logging.TraceId, store)
	)

	recordSvc = record.NewLoggingService(logger, recordSvc)
	bankSvc = bank.NewLoggingService(logger, bankSvc)
	creditCardSvc = creditcard.NewLoggingService(logger, creditCardSvc)
	businessSvc = business.NewLoggingService(logger, businessSvc)
	billSvc = bill.NewLoggingService(logger, billSvc)
	dashboardSvc = dashboard.NewLoggingService(logger, dashboardSvc)
	merchantSvc = merchant.NewLoggingService(logger, merchantSvc)

	userSvc = user.NewLogging(logger, logging.TraceId)(userSvc)
	mpSvc = mp.NewLogging(logger, logging.TraceId)(mpSvc)

	httpLogger := log.With(logger, "component", "http")

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encode.JsonError),
		kithttp.ServerErrorHandler(logging.NewLogErrorHandler(level.Error(logger))),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			guid := request.Header.Get("X-Request-Id")
			token := request.Header.Get("Authorization")

			ctx = context.WithValue(ctx, logging.TraceId, guid)
			ctx = context.WithValue(ctx, "token-context", token)
			return ctx
		}),
		//kithttp.ServerBefore(middleware.TracingServerBefore(tracer)),
	}

	ems := []endpoint.Middleware{
		//middleware.TracingMiddleware(tracer),                                                      // 1
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), 1000)), // 0
	}

	tokenEms := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger, tracer),
	}
	tokenEms = append(tokenEms, ems...)

	r := mux.NewRouter()

	// 以下为系统模块
	// 授权登录模块
	r.PathPrefix("/user").Handler(http.StripPrefix("/user", user.MakeHTTPHandler(userSvc, ems, opts)))

	// 小程序接口
	r.PathPrefix("/mp/api").Handler(http.StripPrefix("/mp/api", mp.MakeHTTPHandler(mpSvc, tokenEms, opts)))

	//mux.Handle("/auth/", auth.MakeHandler(authSvc, httpLogger))
	r.Handle("/record", record.MakeHandler(recordSvc, httpLogger))
	r.Handle("/record/", record.MakeHandler(recordSvc, httpLogger))
	r.Handle("/bank", bank.MakeHandler(bankSvc, httpLogger))
	r.Handle("/bank/", bank.MakeHandler(bankSvc, httpLogger))
	r.Handle("/creditcard", creditcard.MakeHandler(creditCardSvc, httpLogger))
	r.Handle("/creditcard/", creditcard.MakeHandler(creditCardSvc, httpLogger))
	r.Handle("/business", business.MakeHandler(businessSvc, httpLogger))
	r.Handle("/business/", business.MakeHandler(businessSvc, httpLogger))
	r.Handle("/auth/", auth.MakeHandler(authSvc, logger))
	r.Handle("/bill/", bill.MakeHandler(billSvc, logger))
	r.Handle("/bill", bill.MakeHandler(billSvc, logger))
	r.Handle("/dashboard/", dashboard.MakeHandler(dashboardSvc, logger))
	r.Handle("/merchant", merchant.MakeHandler(merchantSvc, logger))

	//r.Handle("/", http.FileServer(http.Dir(cf.GetString("server", "http_static"))))
	//http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir(cf.GetString("server", "http_static")))))

	//http.Handle("/metrics", promhttp.Handler())
	// web页面
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cf.GetString("server", "http_static"))))

	http.Handle("/", accessControl(r, logger))

	{
		cornTab := cron.New()
		billCronJob(cornTab, bill.NewService(logger, store))
		cornTab.Start()
		defer cornTab.Stop()
	}

	initCancelInterrupt()
}

func initCancelInterrupt() {
	errs := make(chan error, 2)
	go func() {
		_ = logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		//errs <- http.ListenAndServe(httpAddr, addCors())
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler, logger log.Logger) http.Handler {
	handlers := make(map[string]string, 3)
	if cf.GetBool("cors", "allow") {
		handlers["Access-Control-Allow-Origin"] = cf.GetString("cors", "origin")
		handlers["Access-Control-Allow-Methods"] = cf.GetString("cors", "methods")
		handlers["Access-Control-Allow-Headers"] = cf.GetString("cors", "headers")
		//reqFun = encode.BeforeRequestFunc(handlers)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range handlers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}
		guid := r.Header.Get(logging.TraceId)
		_ = level.Info(logger).Log(logging.TraceId, guid, "remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}
