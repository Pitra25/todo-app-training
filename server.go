package todo

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, heandler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        heandler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
