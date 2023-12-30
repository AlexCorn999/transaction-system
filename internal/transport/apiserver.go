package transport

import (
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/config"
	"github.com/AlexCorn999/transaction-system/internal/logger"
	"github.com/AlexCorn999/transaction-system/internal/repository"
	"github.com/AlexCorn999/transaction-system/internal/service"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	storage  *repository.Storage
	logger   *log.Logger
	router   *chi.Mux
	config   *config.Config
	invoices *service.Invoices
	withdraw *service.Withdraw
}

func NewAPIServer() *APIServer {
	return &APIServer{
		router: chi.NewRouter(),
		logger: log.New(),
		config: config.NewConfig(),
	}
}

func (s *APIServer) Start() error {
	s.config.ParseFlags()
	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return err
	}

	db, err := s.configureStore()
	if err != nil {
		return err
	}
	s.storage = db
	defer s.storage.Close()

	s.invoices = service.NewInvoices(db)
	s.withdraw = service.NewWithdraw(db, db)

	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Use(logger.WithLogging)
	s.router.Post("/invoice", s.Invoice)
	s.router.Post("/withdraw", s.Withdraw)
}

func (s *APIServer) configureLogger() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureStore() (*repository.Storage, error) {
	db, err := repository.NewStorage(s.config.DataBaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
