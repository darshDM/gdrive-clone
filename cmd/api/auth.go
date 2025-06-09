package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darshDM/gdrive-clone-api/types"
)

func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request) {
	var createRequest types.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := app.userService.CreateNewUser(r.Context(), &createRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "User created successfully")
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest types.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reString, err := app.userService.LoginUser(r.Context(), &loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, reString)
}
