package dbProvider

import (
	"Techiebulter/interview/backend/providers"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type pgClientProvider struct {
	pgClient *sql.DB
	ctx      context.Context
}

func ConnectDB(connectionString string) providers.PgClientProvider {

	var (
		pgClient    *sql.DB
		err         error
		maxAttempts = 3
	)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for i := 0; i < maxAttempts; i++ {
		pgClient, err = sql.Open("postgres", connectionString)
		if err != nil {
			fmt.Println("Unable to create PostgreSQL client:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		err = pgClient.PingContext(ctx)
		if err != nil {
			fmt.Println("Unable to connect to PostgreSQL database:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL client: %v", err)
	} else {
		log.Printf("Successfully connected to PostgreSQL database")
	}

	return &pgClientProvider{
		pgClient: pgClient,
		ctx:      ctx,
	}
}

func (p *pgClientProvider) Ping() error {
	return p.pgClient.Ping()
}

func (p *pgClientProvider) Close() error {
	return p.pgClient.Close()
}

func (p *pgClientProvider) Client() *sql.DB {
	return p.pgClient
}
