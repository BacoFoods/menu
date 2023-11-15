package internal

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Rabbit struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	Q    *amqp.Queue
}

func MustNewRabbitMQ(queueName, host, port string) *Rabbit {
	conn, ch, q, err := NewRabbitMQ(queueName, host, port)
	if err != nil {
		panic(err)
	}

	return &Rabbit{conn, ch, q}
}

func NewRabbitMQ(queueName, host, port string) (*amqp.Connection, *amqp.Channel, *amqp.Queue, error) {
	logrus.Info("connecting to rabbitmq:", host, port)
	conn, err := amqp.Dial("amqp://" + host + ":" + port)
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	logrus.Info("rabbitmq channel acquired")

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err == nil {
		logrus.Info("rabbitmq queue", queueName, "declared")
	}

	return conn, ch, &q, err
}

func (r *Rabbit) PutContent(obj any) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Ch.PublishWithContext(
		ctx,
		"", // TODO: use an per-store exchange
		r.Q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
