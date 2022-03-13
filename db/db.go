package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Url         string
	MaxOpenConn int
	MaxIdleConn int

	// maximum lifetime of database connection in seconds
	ConnMaxLife int
}

func InitDb(cfg Config) (*gorm.DB, error) {
	sqlDb, err := sql.Open("postgres", cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	err = sqlDb.Ping()
	if err != nil {
		return nil, fmt.Errorf("sql.Ping failed: %w", err)
	}

	sqlDb.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDb.SetConnMaxLifetime(
		time.Duration(cfg.ConnMaxLife) * time.Second)

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}

	return gormDb, nil
}
