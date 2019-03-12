// Package server creates a HTTP server which can be used as a registration
// point for webhooks, as well as metrics and everything in between
package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cpssd-students/cheapskate/pkg/config"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

// Error originating from the server itself.
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// InvalidEndpointPrefix occurs when the prefix the handler tried to
	// register is invalid
	InvalidEndpointPrefix = Error("invalid endpoint prefix")
	// EndpointAlreadyUsed indicates that some previous handler already
	// registered on that given endpoint
	EndpointAlreadyUsed = Error("endpoint already used")
)

// Server contains the definitions and implementations of the service
type Server struct {
	mux *http.ServeMux
	// registry of the endpoints
	eps map[string]struct{}
	log zerolog.Logger
	cfg config.Server
}

// New creates a new server with a given configuration
func New(opts ...Option) (*Server, error) {
	s := &Server{
		mux: http.DefaultServeMux,
		eps: make(map[string]struct{}),
		log: zerolog.Nop(),
	}

	// applying options
	for _, o := range opts {
		if err := o.apply(s); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}

	s.register("/_/healthz", s.healthz())
	s.register("/_/metrics", s.metrics())

	// The rest of the handles are registered under the root path, and can be
	// whatever other than /_

	return s, nil
}

// Option adds additional stuff to the server at setup time.
type Option interface {
	apply(*Server) error
}

type optionFunc func(*Server) error

func (of optionFunc) apply(s *Server) error {
	return of(s)
}

// WithLogger allows to specify the zerologger instance to be used.
func WithLogger(l zerolog.Logger) Option {
	return optionFunc(func(s *Server) error {
		s.log = l
		return nil
	})
}

// WithConfig allows to specify the config to be used.
func WithConfig(c config.Server) Option {
	return optionFunc(func(s *Server) error {
		s.cfg = c
		return nil
	})
}

func (s *Server) healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) metrics() http.Handler {
	return promhttp.Handler()
}

func (s *Server) register(ep string, h http.Handler) {
	s.log.Info().
		Str("endpoint", ep).
		Msg("registering endpoint")
	s.eps[ep] = struct{}{}
	s.mux.Handle(ep, h)
}

// Register allows to register a webbook with the http server. If the returned
// response is not nil, it means that the endpoint was not registered.
func (s *Server) Register(ep string, h http.Handler) error {
	if strings.HasPrefix(ep, "/_") || strings.HasPrefix(ep, "_") {
		return InvalidEndpointPrefix
	}

	if _, exists := s.eps[ep]; exists {
		return EndpointAlreadyUsed
	}

	s.register(ep, h)

	return nil
}

// Run the server
func (s *Server) Run() error {
	s.log.Info().Int("port", s.cfg.Port).Msg("starting server")
	return http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), s.mux)
}
