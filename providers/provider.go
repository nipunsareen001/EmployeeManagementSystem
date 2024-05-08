package providers

import (
	"database/sql"
)

// PgClientProvider provides database connection for PostgreSQL.
type PgClientProvider interface {
	// Ping verifies the connection with the database.
	Ping() error

	// Close closes the database connection.
	Close() error

	// Client returns the pointer to the PostgreSQL database client.
	Client() *sql.DB
}
