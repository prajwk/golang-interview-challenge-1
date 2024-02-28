package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gpsinsight/go-interview-challenge/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*http.Server
	router *mux.Router
	logger *logrus.Entry
	db     *sqlx.DB
}

func New(
	cfg config.Config,
	log *logrus.Entry,
	db *sqlx.DB,
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
		db:     db,
	}
	// initialize the routes
	server.Route()
	return server
}
