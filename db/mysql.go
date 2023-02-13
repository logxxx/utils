package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"strings"
	"time"
)

type MysqlClient struct {
	db *sqlx.DB
}

func NewMysqlClient(dataSourceName string,
	maxOpenConns int,
	maxIdleConns int,
	maxLifeTime time.Duration) (*MysqlClient, error) {
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxLifeTime)

	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	client := &MysqlClient{
		db: db,
	}

	return client, nil
}

func (c *MysqlClient) GetDb() *sqlx.DB {
	return c.db
}
