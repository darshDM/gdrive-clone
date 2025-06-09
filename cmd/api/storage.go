package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/darshDM/gdrive-clone-api/types"
)

func (app *application) getRemainingStorageHandler(w http.ResponseWriter, r *http.Request) {
	remainingStorage, err := app.storageService.GetRemainingStorage(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := types.StorageResponse{
		RemainingStorage: remainingStorage,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (app *application) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	const maxSizePerRequest = 10 * 1024 * 1024 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxSizePerRequest)
	if err := r.ParseMultipartForm(maxSizePerRequest); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}
	file, fileInfo, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := app.storageService.UploadFile(r.Context(), file, fileInfo); err != nil {
		if err.Error() == "not enough storage space" {
			http.Error(w, "Not enough storage space", http.StatusBadRequest)
			return
		}

		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := types.FileCreatedResponse{
		Message: "File uploaded successfully",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (app *application) GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
	}

	res, err := app.storageService.GetFiles(r.Context(), limitInt, offsetInt)
	if err != nil {
		http.Error(w, "Failed to get files", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
