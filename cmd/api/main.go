package main

import (
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/db"
	"log"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/env"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@db.adventistasc.orb.local/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
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
