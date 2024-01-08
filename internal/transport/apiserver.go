package transport

import (
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/config"
	"github.com/AlexCorn999/transaction-system/internal/hash"
	"github.com/AlexCorn999/transaction-system/internal/logger"
	"github.com/AlexCorn999/transaction-system/internal/repository"
	"github.com/AlexCorn999/transaction-system/internal/service"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	storage *repository.Storage
	logger  *log.Logger
	router  *chi.Mux
	config  *config.Config
	users   *service.Users
	money   *service.Money
}

func NewAPIServer() *APIServer {
	return &APIServer{
		router: chi.NewRouter(),
		logger: log.New(),
		config: config.NewConfig(),
	}
}

// Start starts and configures the server.
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

	hasher := hash.NewSHA1Hasher("salt")
	s.users = service.NewUsers(db, hasher, []byte("sample secret"), s.config.TokenTTL)

	s.money = service.NewInvoices(db, db)

	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

// configureRouter configures endpoint routing.
func (s *APIServer) configureRouter() {
	s.router.Use(logger.WithLogging)
	s.router.Post("/api/user/register", s.SighUp)
	s.router.Post("/api/user/login", s.SighIn)
	s.router.With(s.authMiddleware).Post("/invoice", s.Invoice)
	s.router.With(s.authMiddleware).Post("/withdraw", s.Withdraw)
	s.router.With(s.authMiddleware).Get("/balance/actual", s.BalanceActual)
	s.router.With(s.authMiddleware).Get("/balance/frozen", s.BalanceFrozen)
}

// configureLogger sets the logger configuration.
func (s *APIServer) configureLogger() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

// configureStore returns an object for working with the database.
func (s *APIServer) configureStore() (*repository.Storage, error) {
	db, err := repository.NewStorage(s.config.DataBaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
