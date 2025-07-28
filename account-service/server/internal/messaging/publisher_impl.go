package messaging

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type AccountCreatedEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *AccountPublisher) PublishAccountCreated(ctx context.Context, id, name, email string) error {
	event := &AccountCreatedEvent{ID: id, Name: name, Email: email}
	body, err := json.Marshal(event)
	if err != nil {
		p.logger.Errorf("failed to marshal account created event: %v", err)
		return err
	}

	p.logger.WithFields(logrus.Fields{"routing_key": accountCreatedRoutingKey, "account_id": id}).Info("Publishing account created event")

	return p.ch.PublishWithContext(
		ctx,
		exchangeName,
		accountCreatedRoutingKey,
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
