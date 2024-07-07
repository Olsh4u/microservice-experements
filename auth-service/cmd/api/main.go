package main

import (
	"database/sql"
	"fmt"
	"github.com/org/olsh4u/auth-service/data"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

func main() {
	log.Println("Starting auth service...")

	// connect to db
	conn := connectDB()
	if conn == nil {
		log.Panic("Could not connect to database")
	}

	//set up config
	app := Config{
		Client: &http.Client{},
	}
	app.setupRepo(conn)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet available...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
