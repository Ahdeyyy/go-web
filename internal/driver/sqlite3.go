package driver

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}
var databasePath = "internal/database/sqlite/"

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSqlite connects to the sqlite database
func ConnectSqlite(dsn string) (*DB, error) {
	var db *sql.DB
	var err error
	if _, err := os.Stat(databasePath + dsn); os.IsNotExist(err) {
		db, err = NewSqliteDB(databasePath + dsn)
		if err != nil {
			panic(err)
		}
	} else {
		db, err = sql.Open("sqlite3", databasePath+dsn)
		if err != nil {
			panic(err)
		}
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = db
	err = testDB(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB tries to ping the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewSqliteDB creates a new sqlite database for the application
func NewSqliteDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")
	db.Exec("PRAGMA busy_timeout = 5000")

	db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, firstname TEXT, lastname TEXT, email TEXT, password TEXT, created_at DATETIME, updated_at DATETIME)")

	return db, nil
}
