package main

import (
	"time"

	"github.com/iamabhishekch/Social/internal/auth"
	"github.com/iamabhishekch/Social/internal/db"
	"github.com/iamabhishekch/Social/internal/env"
	"github.com/iamabhishekch/Social/internal/mailer"
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
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:admin@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdelTime:  env.GetString("DB_MAX_IDLE_CONNS", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("SENDGRID_FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenCofig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "gophersocial",
			},
		},
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

	//mailer
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	// intilizing struct store
	store := store.NewStorage(db)

	//auth
	JWTAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	// ** intilizing struct application **
	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: JWTAuthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
