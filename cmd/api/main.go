package main

import (
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/db"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/env"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store"
	"log"
)

const version = "0.0.1"

//	@title			Swagger Example API
//	@description	Api for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @SecurityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@db.adventistasc.orb.local/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "dev"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Printf("connected to database")

	appStore := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  appStore,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
