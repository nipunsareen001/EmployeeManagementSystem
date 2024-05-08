package utils

import (
	"Techiebulter/interview/backend/models"
	"os"
)

// GetPGSQLConnectionString gets  psqlDB URL from the environment variables
func GetPGSQLConnectionString() string {
	return os.Getenv(string(models.PGSQL_URL))
}

// GetGINPORTString gets FIBER Listener port from the environment variables
func GetFIBERPORTString() string {
	return os.Getenv(string(models.FIBER_PORT))
}
