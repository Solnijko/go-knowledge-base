package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	db, err := goose.OpenDBWithDriver("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	goose.SetBaseFS(os.DirFS("migrations"))
	err = goose.Up(db, ".")
	if err != nil {
		log.Fatal(err)
	}
}
