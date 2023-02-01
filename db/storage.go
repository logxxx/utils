package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	client *MysqlClient
)

func NewClient(cfg DBConfig) (*MysqlClient, error) {
	c, err := NewMysqlClient(
		getDSN(cfg),
		60,
		30,
		5*time.Second)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Init(cfg DBConfig) error {
	var err error
	client, err = NewClient(cfg)
	if err != nil {
		return err
	}
	return nil
}

func getDSN(conf DBConfig) string {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	return dsn
}

func GetDb() *sqlx.DB {
	return client.GetDb()
}
