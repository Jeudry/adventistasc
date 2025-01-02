package main

import (
	"log"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/env"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store: store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
