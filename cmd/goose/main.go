package main

import (
	"database/sql"
	"embed"
	"flag"
	"log"
	"os"

	_ "weight-tracker/cmd/goose/migrations"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
//go:embed migrations/*.go
var embedMigrations embed.FS

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	down   = flags.Bool("down", false, "Migrate the database down")
	downone = flags.Bool("downone", false, "Migrate the database down one step")
	upone = flags.Bool("upone", false, "Migrate the database up one step")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}

	// setup database
	dburl := os.Getenv("BLUEPRINT_DB_URL")
	if dburl == "" {
		panic("BLUEPRINT_DB_URL is empty")
	}
	db, err := sql.Open("sqlite3", dburl)

	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}


	if *down {
		if err := goose.DownTo(db, "migrations", 0); err != nil {
			panic(err)
		}
		return
	}

	if *downone {
		if err := goose.Down(db, "migrations"); err != nil {
			panic(err)
		}
		return
	}

	if *upone {
		if err := goose.UpByOne(db, "migrations"); err != nil {
			panic(err)
		}
		return
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
