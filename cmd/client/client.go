/**
 * @Time: 2019-08-18 16:51
 * @Author: solacowa@gmail.com
 * @File: client
 * @Software: GoLand
 */

package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/nsini/cardbill/src/repository/types"
	"os"
)

var (
	fs         = flag.NewFlagSet("cardbill", flag.ExitOnError)
	configFile = fs.String("config-file", "app.cfg", "server config file")
)

var logger log.Logger

func main() {
	logger = log.NewLogfmtLogger(log.StdlibWriter{})
	logger = log.With(logger, "caller", log.DefaultCaller)

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
	db, err := mysqlclient.NewMysql(dbUrl, cf.GetBool(config.SectionServer, "debug"))
	if err != nil {
		_ = level.Error(logger).Log("db", "connect", "err", err)
		return
	}

	_ = logger.Log("create", "table", "Business", db.CreateTable(types.Business{}).Error)
	_ = logger.Log("create", "table", "Bank", db.CreateTable(types.Bank{}).Error)
	_ = logger.Log("create", "table", "CreditCard", db.CreateTable(types.CreditCard{}).Error)
	_ = logger.Log("create", "table", "ExpensesRecord", db.CreateTable(types.ExpensesRecord{}).Error)
	_ = logger.Log("create", "table", "Rate", db.CreateTable(types.Rate{}).Error)
	_ = logger.Log("create", "table", "User", db.CreateTable(types.User{}).Error)
	_ = logger.Log("create", "table", "Merchant", db.CreateTable(types.Merchant{}).Error)
	_ = logger.Log("create", "table", "Bill", db.CreateTable(types.Bill{}).Error)
}
