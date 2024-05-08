package dbHelperProvider

import (
	"Techiebulter/interview/backend/providers"
	"database/sql"

	_ "github.com/lib/pq"
)

type DBHelper struct {
	pgClient *sql.DB
}

func NewDBHelper(pgClient *sql.DB) providers.DbHelperProvider {
	return &DBHelper{
		pgClient: pgClient,
	}
}
