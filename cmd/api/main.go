package main

import (
	"log"

	"github.com/iamabhishekch/Social/internal/db"
	"github.com/iamabhishekch/Social/internal/env"
	"github.com/iamabhishekch/Social/internal/store"
)

const version = "0.0.1"

func main() {

	// intilizing struct config
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:admin@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdelTime:  env.GetString("DB_MAX_IDLE_CONNS", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// db connection establish
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdelTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("DB connection established")

	// intilizing struct store
	store := store.NewStorage(db)

	// ** intilizing struct application **
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
