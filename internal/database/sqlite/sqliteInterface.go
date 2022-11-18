package sqlite

import (
	"database/sql"

	"github.com/Ahdeyyy/go-web/internal/config"
	"github.com/Ahdeyyy/go-web/internal/database"
)

type sqliteDBinterface struct {
	App *config.Config
	DB  *sql.DB
}

func NewSqliteInterface(conn *sql.DB, a *config.Config) database.DBInterface {
	return &sqliteDBinterface{
		App: a,
		DB:  conn,
	}
}
