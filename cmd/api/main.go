package main

import (
	"github.com/darshDM/gdrive-clone-api/internal/storage"
	"github.com/darshDM/gdrive-clone-api/internal/store"
	"github.com/darshDM/gdrive-clone-api/internal/user"
)

func main() {
	// absPath, err := filepath.Abs("cmd/migrations/")
	// if err != nil {
	// 	log.Println("Error getting absolute path:", err)
	// }
	databaseConfig := &store.DatabaseConfig{
		DatabaseName:    "db/gdrive-db",
		MigrationFolder: "./cmd/migrations/",
	}
	sqlite := store.NewSQLiteStore(databaseConfig)
	userService := user.NewUserService(sqlite)
	storageService := storage.NewStorageService(sqlite)

	app := &application{
		userService:    userService,
		storageService: storageService,
	}

	r := app.Mount()
	app.Run(r)

}
