package main

import (
	"github.com/iamabhishekch/Social/internal/db"
	"github.com/iamabhishekch/Social/internal/env"
	"github.com/iamabhishekch/Social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Swagger GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host						petstore.swagger.io
//	@BasePath					/v1
//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
// description

func main() {

	// intilizing struct config
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:admin@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdelTime:  env.GetString("DB_MAX_IDLE_CONNS", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	// sync flashed any log entires
	defer logger.Sync()

	// db connection establish
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdelTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB connection established")

	// intilizing struct store
	store := store.NewStorage(db)

	// ** intilizing struct application **
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
