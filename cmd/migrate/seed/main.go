package main

import (
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/db"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/env"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@db.adventistasc.orb.local/socialnetwork?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)
}
