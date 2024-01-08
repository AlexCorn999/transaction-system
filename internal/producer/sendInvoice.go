package producer

import (
	"io"
	"net/http"
	"strconv"

	"github.com/AlexCorn999/transaction-system/internal/domain"
	"github.com/AlexCorn999/transaction-system/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (s *APIServer) SendInvoice(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("sendInvoice", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		logger.LogError("sendInvoice", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.rabbitmq.Ch.PublishWithContext(r.Context(),
		"",
		"transactions",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			MessageId:   strconv.Itoa(int(userID)),
		})

	s.logger.Info("message has been forwarded")
	if err != nil {
		logger.LogError("sendInvoice", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
