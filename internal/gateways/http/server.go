package http

import (
	"context"
	"fmt"
	"homework/internal/usecase"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	port    uint16
	router  *gin.Engine
	handler *WebSocketHandler
}

type UseCases struct {
	Event  *usecase.Event
	Sensor *usecase.Sensor
	User   *usecase.User
}

func NewServer(useCases UseCases, options ...func(*Server)) *Server {
	r := gin.Default()
	handler := NewWebSocketHandler(useCases)
	setupRouter(r, useCases, handler)

	s := &Server{router: r, host: "localhost", port: 8080, handler: handler}
	for _, o := range options {
		o(s)
	}

	return s
}

func WithHost(host string) func(*Server) {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.router,
	}

	var serverErr error
	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		serverErr = srv.ListenAndServe()
	}()
	wg.Add(1)

	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error while closing server: %v", err)
	}
	if err := s.handler.Shutdown(); err != nil {
		log.Printf("error while closing web-socket: %v", err)
	}

	wg.Wait()

	return serverErr
}
