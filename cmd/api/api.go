package main

import (
	"log"
	"net/http"

	"github.com/darshDM/gdrive-clone-api/internal/storage"
	"github.com/darshDM/gdrive-clone-api/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	userService    *user.UserService
	storageService *storage.StorageService
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (app *application) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/users", func(r chi.Router) {
			r.Post("/signup", app.signUpHandler)
			r.Post("/login", app.loginHandler)

		})
		r.Route("/storage", func(r chi.Router) {
			r.Use(app.Authenticate)
			r.Get("/files", app.GetFilesHandler)
			r.Get("/remaining", app.getRemainingStorageHandler)
			r.Post("/upload", app.UploadFileHandler)
		})
	})
	return r
}

func (app *application) Run(h http.Handler) {
	log.Println("Starting server on :8000")
	err := http.ListenAndServe(":8000", h)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc5OTgzNjksInVzZXJuYW1lIjoiZGFyc2gifQ.1IoN4RN8mvNU-r-mQ-klz5xv_VPirZqMrR7e_2V4hZs
