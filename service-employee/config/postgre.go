package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	// Set your PostgreSQL connection parameters
	// Example URI: "user=username password=password dbname=service_employee sslmode=disable"
	uri := os.Getenv("POSTGRESQL")

	// Connect to the PostgreSQL database
	conn, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}

	// Ping the database to confirm a successful connection
	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to PostgreSQL!")

	db = conn
}

func NewPostgresContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func GetPostgresDB() *sql.DB {
	return db
}
