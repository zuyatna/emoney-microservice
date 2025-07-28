package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const exchangeName = "emoney_exchange"
const accountCreatedRoutingKey = "account.created"

type AccountPublisher struct {
	ch     *amqp.Channel
	logger *logrus.Logger
}

func NewAccountPublisher(conn *amqp.Connection, logger *logrus.Logger) (*AccountPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic", // Exchange type
		true,    // Durable
		false,   // Auto-deleted
		false,   // Internal
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		err := ch.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &AccountPublisher{
		ch:     ch,
		logger: logger,
	}, nil
}
