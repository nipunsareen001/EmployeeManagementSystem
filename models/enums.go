package models

type DatabaseURL string
type PORT string

const (
	PGSQL_URL  DatabaseURL = "PGSQL_URL"
	FIBER_PORT PORT        = "FIBER_PORT"
)
