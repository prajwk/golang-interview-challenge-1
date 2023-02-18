package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gpsinsight/go-interview-challenge/internal/config"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*http.Server
	router *mux.Router
	logger *logrus.Entry
}

func New(
	cfg config.Config,
	log *logrus.Entry,
) *Server {
	timeout := 60 * time.Second
	router := mux.NewRouter()
	router.StrictSlash(true)

	server := &Server{
		Server: &http.Server{
			Handler:           router,
			Addr:              ":80",
			ReadTimeout:       timeout,
			WriteTimeout:      timeout,
			ReadHeaderTimeout: timeout,
			IdleTimeout:       timeout,
			MaxHeaderBytes:    65536,
		},
		router: router,
		logger: log,
	}

	return server
}
