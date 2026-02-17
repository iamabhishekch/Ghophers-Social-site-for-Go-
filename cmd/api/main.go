package main

import (
	"log"

	"github.com/iamabhishekch/Social/internal/db"
	"github.com/iamabhishekch/Social/internal/env"
	"github.com/iamabhishekch/Social/internal/store"
)

func main() {

	// intilizing struct config
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:admin@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdelTime:  env.GetString("DB_MAX_IDLE_CONNS", "15m"),
		},
	}


	// db connection establish
	db, err := db.New(
		cfg.db.addr, 
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdelTime,
	)
	if err!= nil{
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
