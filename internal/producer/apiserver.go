package producer

import (
	"net/http"

	"github.com/AlexCorn999/transaction-system/internal/config"
	"github.com/AlexCorn999/transaction-system/internal/hash"
	"github.com/AlexCorn999/transaction-system/internal/logger"
	"github.com/AlexCorn999/transaction-system/internal/rabbitmq"
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
	rabbitmq *rabbitmq.RabbitMQ
	config   *config.Config
	users    *service.Users
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

	if err := s.configureRabbitMQ(); err != nil {
		return err
	}
	defer s.rabbitmq.Close()

	db, err := s.configureStore()
	if err != nil {
		return err
	}
	s.storage = db
	defer s.storage.Close()

	hasher := hash.NewSHA1Hasher("salt")
	s.users = service.NewUsers(db, hasher, []byte("sample secret"), s.config.TokenTTL)

	s.logger.Info("starting sender")

	return http.ListenAndServe(":3000", s.router)
}

// configureRouter configures endpoint routing.
func (s *APIServer) configureRouter() {
	s.router.Use(logger.WithLogging)
	s.router.Post("/api/user/register", s.SighUp)
	s.router.Post("/api/user/login", s.SighIn)
	s.router.With(s.authMiddleware).Post("/invoice", s.SendInvoice)
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

// configureRabbitMQ connects to rabbitmq.
func (s *APIServer) configureRabbitMQ() error {
	rabbitmq, err := rabbitmq.NewRabbitMQ()
	if err != nil {
		return err
	}
	s.rabbitmq = rabbitmq
	return nil
}
