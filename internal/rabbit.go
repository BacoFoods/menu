package internal

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
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
	conn, err := amqp.Dial("amqp://" + host + ":" + port)
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	return conn, ch, &q, err
}

func (r *Rabbit) PutContent(obj any) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return r.Ch.Publish(
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
