/**
 * @Time: 2019-08-18 09:48
 * @Author: solacowa@gmail.com
 * @File: main
 * @Software: GoLand
 */

package service

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/pkg/auth"
	"github.com/nsini/cardbill/src/pkg/bank"
	"github.com/nsini/cardbill/src/pkg/bill"
	"github.com/nsini/cardbill/src/pkg/business"
	"github.com/nsini/cardbill/src/pkg/creditcard"
	"github.com/nsini/cardbill/src/pkg/dashboard"
	"github.com/nsini/cardbill/src/pkg/merchant"
	"github.com/nsini/cardbill/src/pkg/record"
	"github.com/nsini/cardbill/src/pkg/user"
	"github.com/nsini/cardbill/src/repository"
	"github.com/robfig/cron"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var logger log.Logger

var (
	fs         = flag.NewFlagSet("cardbill", flag.ExitOnError)
	httpAddr   = fs.String("http-addr", ":8080", "HTTP listen address")
	configFile = fs.String("config-file", "app.cfg", "server config file")

	db *gorm.DB
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

	cf, err := config.NewConfig(*configFile)
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

	store := repository.NewRepository(db)

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
	)

	recordSvc = record.NewLoggingService(logger, recordSvc)
	bankSvc = bank.NewLoggingService(logger, bankSvc)
	creditCardSvc = creditcard.NewLoggingService(logger, creditCardSvc)
	userSvc = user.NewLoggingService(logger, userSvc)
	businessSvc = business.NewLoggingService(logger, businessSvc)
	billSvc = bill.NewLoggingService(logger, billSvc)
	dashboardSvc = dashboard.NewLoggingService(logger, dashboardSvc)
	merchantSvc = merchant.NewLoggingService(logger, merchantSvc)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	//mux.Handle("/auth/", auth.MakeHandler(authSvc, httpLogger))
	mux.Handle("/record", record.MakeHandler(recordSvc, httpLogger))
	mux.Handle("/record/", record.MakeHandler(recordSvc, httpLogger))
	mux.Handle("/bank", bank.MakeHandler(bankSvc, httpLogger))
	mux.Handle("/bank/", bank.MakeHandler(bankSvc, httpLogger))
	mux.Handle("/creditcard", creditcard.MakeHandler(creditCardSvc, httpLogger))
	mux.Handle("/creditcard/", creditcard.MakeHandler(creditCardSvc, httpLogger))
	mux.Handle("/user/", user.MakeHandler(userSvc, httpLogger))
	mux.Handle("/business", business.MakeHandler(businessSvc, httpLogger))
	mux.Handle("/business/", business.MakeHandler(businessSvc, httpLogger))
	mux.Handle("/auth/", auth.MakeHandler(authSvc, logger))
	mux.Handle("/bill/", bill.MakeHandler(billSvc, logger))
	mux.Handle("/bill", bill.MakeHandler(billSvc, logger))
	mux.Handle("/dashboard/", dashboard.MakeHandler(dashboardSvc, logger))
	mux.Handle("/merchant", merchant.MakeHandler(merchantSvc, logger))

	mux.Handle("/", http.FileServer(http.Dir(cf.GetString("server", "http_static"))))
	//http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir(cf.GetString("server", "http_static")))))

	//http.Handle("/metrics", promhttp.Handler())

	handlers := make(map[string]string, 3)
	if cf.GetBool("cors", "allow") {
		handlers["Access-Control-Allow-Origin"] = cf.GetString("cors", "origin")
		handlers["Access-Control-Allow-Methods"] = cf.GetString("cors", "methods")
		handlers["Access-Control-Allow-Headers"] = cf.GetString("cors", "headers")
	}
	http.Handle("/", accessControl(mux, logger, handlers))

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

func accessControl(h http.Handler, logger log.Logger, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		//requestId := r.Header.Get("X-Request-Id")
		_ = logger.Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})
}
