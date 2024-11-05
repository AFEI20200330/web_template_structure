package mysql

import (
	"fmt"
	"web_template/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	// init mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	
	// connect mysql,sqlx is a high-performance mysql driver for Go
	db, err := sqlx.Connect("mysql", dsn) //
	if err != nil {
		zap.L().Error("connect mysql error", zap.Error(err))
		return err
	}
	// set max open connections
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	// set max idle connections
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return nil
}

func Close(){
	if err := db.Close(); err != nil {
		zap.L().Error("close mysql error", zap.Error(err))
	}
}
