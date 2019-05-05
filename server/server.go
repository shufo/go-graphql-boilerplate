package server

import (
	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Server represents main server structure
type Server struct {
	router *chi.Mux
	config Config
}

type Config struct {
	Logging bool
}

// NewServer returns server with initialized router
func NewServer(c Config) *Server {
	return &Server{
		router: chi.NewRouter(),
		config: c,
	}
}

func init() {
}
