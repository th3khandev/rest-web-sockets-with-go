package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"github.com/th3khan/rest-web-sockets-with-go/models"
	"github.com/th3khan/rest-web-sockets-with-go/repositories"
	"github.com/th3khan/rest-web-sockets-with-go/server"
)

type InsertPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostResponse struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest InsertPostRequest
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post := models.Post{
				ID:      id.String(),
				UserID:  claims.UserID,
				Title:   postRequest.Title,
				Content: postRequest.Content,
			}
			err = repositories.InsertPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			postResponse := PostResponse{
				ID:      post.ID,
				UserID:  post.UserID,
				Title:   post.Title,
				Content: post.Content,
			}
			json.NewEncoder(w).Encode(postResponse)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	}
}
