package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/th3khan/rest-web-sockets-with-go/database"
	"github.com/th3khan/rest-web-sockets-with-go/repositories"
)

type Config struct {
	Port        string
	JWTSecret   string
	DataBaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("Secret Key is Required")
	}
	if config.DataBaseUrl == "" {
		return nil, errors.New("Database url is required")
	}
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	repo, err := database.NewMySQLRepository(b.config.DataBaseUrl)
	if err != nil {
		log.Fatal("Error", err)
	}

	repositories.SetRepository(repo)

	log.Println("Starting server on Port", b.config.Port)

	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
