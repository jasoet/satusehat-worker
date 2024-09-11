package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	"strings"
	"time"
)

type DatabaseType string

const (
	Mysql      DatabaseType = "MYSQL"
	Postgresql DatabaseType = "POSTGRES"
	MSSQL      DatabaseType = "MSSQL"
)

type ConnectionConfig struct {
	DbType       DatabaseType  `yaml:"db_type" validate:"required,oneof=MYSQL POSTGRES MSSQL" mapstructure:"db_type"`
	Host         string        `yaml:"host" validate:"required,min=1" mapstructure:"host"`
	Port         int           `yaml:"port" mapstructure:"port"`
	Username     string        `yaml:"username" validate:"required,min=1" mapstructure:"username"`
	Password     string        `yaml:"password" mapstructure:"password"`
	DbName       string        `yaml:"db_name" validate:"required,min=1" mapstructure:"db_name"`
	Timeout      time.Duration `yaml:"timeout" mapstructure:"timeout" validate:"min=3s"`
	MaxIdleConns int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns" validate:"min=1"`
	MaxOpenConns int           `yaml:"max_open_conns" mapstructure:"max_open_conns" validate:"min=2"`
}

func (c *ConnectionConfig) Dsn() string {
	timeoutString := fmt.Sprintf("%ds", c.Timeout/time.Second)

	var dsn string
	switch c.DbType {
	case Mysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=%s", c.Username, c.Password, c.Host, c.Port, c.DbName, timeoutString)
	case Postgresql:
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable connect_timeout=%s", c.Username, c.Password, c.Host, c.Port, c.DbName, timeoutString)
	case MSSQL:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&connectTimeout=%s", c.Username, c.Password, c.Host, c.Port, c.DbName, timeoutString)
	}

	return dsn
}

func (c *ConnectionConfig) Pool() (*sqlx.DB, error) {
	if c.Dsn() == "" {
		return nil, fmt.Errorf("dsn is empty")
	}

	dbType := string(c.DbType)
	db, err := sqlx.Open(strings.ToLower(dbType), c.Dsn())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
