package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Url         string
	MaxOpenConn int
	MaxIdleConn int
	ConnMaxLife time.Duration
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
	sqlDb.SetConnMaxLifetime(cfg.ConnMaxLife)

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}

	migrateDb(sqlDb)

	return gormDb, nil
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

//go:embed migrations/*.sql
var migrations embed.FS

func migrateDb(db *sql.DB) {
	// Create table to store executed migrations
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS _migrations (
		id int PRIMARY KEY,
		timestamp timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`)
	checkErr(err)

	count := 0

	for true {
		// Start transaction
		tx, err := db.Begin()
		checkErr(err)
		defer func() {
			tx.Rollback()
		}()

		// Read next migration from database, if not found use 1
		var next int
		checkErr(tx.QueryRow(`SELECT COALESCE(MAX(id),0)+1 FROM _migrations`).Scan(&next))

		// Try to read and execute next migration
		file := fmt.Sprintf("migrations/%06d", next) + ".sql"
		if bytes, err := migrations.ReadFile(file); err == nil {
			log.Printf("Running migration %06d...\n", next)
			_, err = tx.Exec(string(bytes))
			checkErr(err)
		} else {
			// File does not exist, no migration to run
			if count == 0 {
				log.Println("Database already up to date.")
			} else {
				log.Println("Database successfully migrated.")
			}

			return
		}

		// Log execution of migration
		_, err = tx.Exec(`INSERT INTO _migrations(id) VALUES ($1)`, next)
		checkErr(err)

		// Commit
		checkErr(tx.Commit())

		count++
	}

	panic("Application should never reach this point")
}
