package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

var DB *sql.DB

func DatabaseConnection(ctx context.Context) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	sqltrace.Register("postgres", &pq.Driver{},
		sqltrace.WithServiceName("task-manager-db"),
	)

	var err error
	DB, err = sqltrace.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	if err = DB.PingContext(ctx); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	log.Println("Connected to DB with Datadog tracing")
}
