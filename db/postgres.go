package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Run() {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://root:pass@localhost:8035/tasker")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// caller must close conn
	defer pool.Close()

	var greeting string
	err = pool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	err = Migrate(pool)
	if err != nil {
		log.Println("Error migrating", err.Error())
	}

	fmt.Println(greeting)
}

func Migrate(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS Test (
			id serial PRIMARY KEY,
			name varchar(20) NOT NULL,
			created date
			);
		`)
	return err
}

type Thing struct {
	Name    string
	Created time.Time
}
