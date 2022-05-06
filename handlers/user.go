package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/segmentio/ksuid"
	"github.com/th3khan/rest-web-sockets-with-go/models"
	"github.com/th3khan/rest-web-sockets-with-go/repositories"
	"github.com/th3khan/rest-web-sockets-with-go/server"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request SignUpRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := models.User{
			ID:       id.String(),
			Email:    request.Email,
			Password: string(hashedPassword),
		}

		err = repositories.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SignUpResponse{
			ID:    user.ID,
			Email: user.Email,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
