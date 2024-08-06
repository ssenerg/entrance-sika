package api

import (
	"entrance/database"
	"log"

	"github.com/labstack/echo/v4"
)

type Server struct {
	port string
	db   *database.Database
	echo *echo.Echo
}

func NewServer(port string, dbConfig *database.Config) (*Server, error) {
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		port: port,
		db:   db,
		echo: echo.New(),
	}, nil
}

func (s *Server) RegisterHandlers() {
	s.echo.GET("/users/:id", s.GetUserByID)
}

func (s *Server) Start() {
	log.Printf("Server starting on port %s", s.port)
	err := s.echo.Start(s.port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
