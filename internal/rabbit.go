package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Rabbit struct {
	host      string
	port      string
	queueName string

	Conn *amqp.Connection
	Ch   *amqp.Channel
	Q    *amqp.Queue

	isReconnecting bool
	mu             *sync.Mutex
}

func MustNewRabbitMQ(queueName, host, port string) *Rabbit {
	r, err := NewRabbitMQ(queueName, host, port)
	if err != nil {
		panic(err)
	}

	return r
}

func (r *Rabbit) reconnect() error {
	r.isReconnecting = true

	r.mu.Lock()
	defer r.mu.Unlock()

	// close channel
	if r.Ch != nil {
		r.Ch.Close()
	}

	// close connection
	if r.Conn != nil {
		r.Conn.Close()
	}

	conn, err := amqp.Dial("amqp://" + r.host + ":" + r.port)
	if err != nil {
		return err
	}
	r.Conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	r.Ch = ch

	notifyConnCloseCh := r.Conn.NotifyClose(make(chan *amqp.Error))
	notifyChanCloseCh := r.Ch.NotifyClose(make(chan *amqp.Error))

	go func() {
		for notifyChanCloseCh != nil || notifyConnCloseCh != nil {
			select {
			case err, ok := <-notifyConnCloseCh:
				if !ok {
					notifyConnCloseCh = nil
				} else {
					fmt.Printf("connection closed, error %s ... reconnecting", err)
					if err := r.reconnect(); err != nil {
						logrus.Error("error reconnecting to rabbitmq:", err)
					}
				}
			case err, ok := <-notifyChanCloseCh:
				if !ok {
					notifyChanCloseCh = nil
				} else {
					fmt.Printf("channel closed, error %s ... reconnecting", err)
					if err := r.reconnect(); err != nil {
						logrus.Error("error reconnecting to rabbitmq:", err)
					}
				}
			}
		}
	}()

	logrus.Info("rabbitmq channel acquired")

	q, err := ch.QueueDeclare(
		r.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	r.Q = &q

	if err == nil {
		logrus.Info("rabbitmq queue", r.queueName, "declared")
	}

	r.isReconnecting = false

	return err
}

func (r *Rabbit) ensureConnection() error {
	if !r.Conn.IsClosed() && !r.Ch.IsClosed() {
		// connection is stable
		return nil
	}

	// is already re-connecting
	if r.isReconnecting {
		// wait for connection to be re-established
		r.mu.Lock()
		r.mu.Unlock()
		return nil
	}

	return r.reconnect()
}

func NewRabbitMQ(queueName, host, port string) (*Rabbit, error) {
	logrus.Info("connecting to rabbitmq:", host, port)

	r := &Rabbit{
		host:           host,
		port:           port,
		queueName:      queueName,
		isReconnecting: false,
		mu:             &sync.Mutex{},
	}

	err := r.reconnect()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rabbit) PutContent(obj any) error {
	if err := r.ensureConnection(); err != nil {
		return err
	}

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
