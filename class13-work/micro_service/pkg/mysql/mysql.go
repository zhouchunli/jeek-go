package mysql

import (
	"log"
	"micro_service/pkg/mysql/sqlt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Mysql
type Conf struct {
	Address     string
	MaxOpenConn []int
	MaxIdleConn []int
	MaxLifetime []time.Duration
}

func New(cfg *Conf) *sqlt.DB {
	mysqlDB, err := sqlt.Open("mysql", cfg.Address)

	if err != nil {
		panic(err)
	}
	err = mysqlDB.Ping()
	if err != nil {
		panic(err)
	}
	mysqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	mysqlDB.SetMaxOpenConnections(cfg.MaxOpenConn)
	for i, v := range cfg.MaxLifetime {
		cfg.MaxLifetime[i] = v * time.Second
	}
	mysqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
	log.Println("Mysql is connected!!!")
	return mysqlDB
}
