package transport

import (
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/config"
	"github.com/AlexCorn999/transaction-system/internal/logger"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	//store  store.Store
	logger *log.Logger
	router *chi.Mux
	config *config.Config
}

func NewAPIServer() *APIServer {
	return &APIServer{
		router: chi.NewRouter(),
		logger: log.New(),
		//	store:  store,
		config: config.NewConfig(),
	}
}

func (s *APIServer) Start() error {
	s.config.ParseFlags()
	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Use(logger.WithLogging)
	s.router.Get("/invoice", s.Invoice)
}

func (s *APIServer) configureLogger() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}
