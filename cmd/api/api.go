package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Post("/{postId}/comments", app.createPostCommentHandler)
				r.Put("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			//r.Post("/", app.createUserHandler)
			//r.Get("/", app.getUsersHandler)

			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/toggle_follow", app.toggleFollowUserHandler)
				//r.Put("/", app.updateUserHandler)
				//r.Delete("/", app.deleteUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
