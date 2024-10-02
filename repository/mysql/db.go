package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	UserName string `koanf:"username"`
	Password string `koanf:"password"`
	Port     int    `koanf:"port"`
	Host     string `koanf:"host"`
	DBName   string `koanf:"db_name"`
}

type MySQLDB struct {
	config Config
	db     *sql.DB
}

func (m *MySQLDB) Conn() *sql.DB {
	return m.db
}

func New(cfg Config) *MySQLDB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{config: cfg, db: db}
}
