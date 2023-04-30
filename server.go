package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// ServerOpts specifies the options with which to create an App.
type ServerOpts struct {
	Host string
	Port int

	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string
}

// Validate ensures that the options provided are sane.
func (o *ServerOpts) Validate() error {
	switch {
	case o.Port > 65535:
		return fmt.Errorf("supplied Port higher than max port (65535)")
	case o.TLSEnabled:
		if o.TLSCertFile == "" {
			return fmt.Errorf("TLSEnabled was set, but TLSCertFile was empty")
		} else if o.TLSKeyFile == "" {
			return fmt.Errorf("TLSEnabled was set, but TLSKeyFile was empty")
		}
	}

	return nil
}

// Server holds the configuration of the application.
type Server struct {
	*ServerOpts
	Server *http.Server
	Router *httprouter.Router

	searchEngines SearchEngines
	initialized   bool
}

// New returns a new App.
func NewServer(opts *ServerOpts) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	// TLS
	addr := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	tlsConf := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	// router := http.NewServeMux()
	router := httprouter.New()

	// General HTTP server
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		TLSConfig:    tlsConf,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	app := &Server{
		ServerOpts: opts,

		Server: server,
		Router: router,
	}

	if err := app.setupHandlers(); err != nil {
		return nil, fmt.Errorf("failed to set up handlers: %w", err)
	}

	app.initialized = true

	return app, nil
}

func (s *Server) setupHandlers() error {
	s.Router.GET("/", s.Index)
	s.Router.GET("/osxml/:name", s.OSXML)

	return nil
}

// Start kicks off the Server.
func (s *Server) Start() error {
	var err error

	if s.TLSEnabled {
		err = s.Server.ListenAndServeTLS(s.TLSCertFile, s.TLSKeyFile)
	} else {
		err = s.Server.ListenAndServe()
	}

	return fmt.Errorf("server error: %w", err)
}

// LoadSearchEngines validates each of the SearchEngines and sets them on the Server, returning an error if unsuccessful.
func (s *Server) LoadSearchEngines(se SearchEngines) error {
	if err := se.Validate(); err != nil {
		return err
	}

	s.searchEngines = se

	return nil
}
