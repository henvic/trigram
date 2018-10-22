package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/henvic/trigram"
)

// Params for the server.
type Params struct {
	Address string
}

// Run the server.
func Run(ctx context.Context, p Params) error {
	var mux = http.NewServeMux()

	var s = NewServer()

	mux.HandleFunc("/learn", s.learnHandler)
	mux.HandleFunc("/generate", s.generateHandler)

	log.Println("Exposing HTTP server on", p.Address)

	var server = &http.Server{
		Addr:    p.Address,
		Handler: mux,
	}

	ec := make(chan error, 1)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil && err != context.Canceled {
			ec <- err
		}
	}()

	go func() {
		ec <- server.ListenAndServe()
	}()

	return <-ec
}

// NewServer creates a new server instance.
func NewServer() *Server {
	return &Server{
		store: trigram.NewStore(),
	}
}

// Server for handling Trigram requests.
type Server struct {
	store *trigram.Store
}

func (s *Server) learnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		statusHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		statusHandler(w, r, http.StatusNotAcceptable)
		return
	}

	if err := s.store.Learn(r.Context(), r.Body); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		statusHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	generated, err := s.store.Generate()

	switch err {
	case nil:
		_, _ = fmt.Fprintf(w, "%s\n", generated)
	case trigram.ErrTooShort:
		errorHandler(w, r, http.StatusNotFound, err)
	default:
		errorHandler(w, r, http.StatusInternalServerError, err) // ignoring logging errors for this version
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "%v %v\n", status, http.StatusText(status))
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "%v %v\n%v\n", status, http.StatusText(status), err)
}
