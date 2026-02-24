package main

import (
	"log"

	"github.com/iamabhishekch/Social/internal/db"
	"github.com/iamabhishekch/Social/internal/env"
	"github.com/iamabhishekch/Social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:admin@localhost:5433/social?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store)
}
