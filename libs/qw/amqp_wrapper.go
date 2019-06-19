package qw

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Q represents the queue
type Q struct {
	Endpoint string
	ID       string

	server  *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue

	brokerError chan *amqp.Error
}

// Declare declares a new queue
func (q *Q) initialize() error {
	var err error

	q.brokerError = make(chan *amqp.Error)

	if q.server, err = amqp.Dial(q.Endpoint); err != nil {
		log.Println("Error connecting to RabbitMQ")
		return err
	}

	q.server.NotifyClose(q.brokerError)

	if q.channel, err = q.server.Channel(); err != nil {
		log.Println("Connected creating channel")
		return err
	}

	if q.queue, err = q.channel.QueueDeclare(q.ID, true, false, false, false, nil); err != nil {
		log.Println("Error declaring queue")
		return err
	}

	if err = q.channel.Qos(1, 0, false); err != nil {
		log.Println("Failed to set QoS")
		return err
	}

	return err
}

// Create instanciates a new queue with the given parameters
func Create(endpoint string, id string) *Q {
	instance := &Q{
		Endpoint: endpoint,
		ID:       id,
	}

	return instance
}

// Start starts the queue
func (q *Q) Start() chan bool {
	ready := make(chan bool)

	go func() {
		for {
			if err := q.initialize(); err == nil {
				ready <- true
				q.handleError()
			}

			log.Printf("Connection to RabbitMQ lost, trying to reconnect ...")
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	return ready
}

// Consume returns the messages channel to be readed
func (q *Q) Consume() (<-chan amqp.Delivery, error) {
	return q.channel.Consume(q.queue.Name, "", false, false, false, false, nil)
}

func (q *Q) handleError() {
	err := <-q.brokerError
	log.Printf("Broker error: %v", err)
	q.channel.Close()
	q.server.Close()
}

// Close closes the channel and the server connection
func (q *Q) Close() {
	q.channel.Close()
	q.server.Close()
}

// Write stacks the given message in the queue
func (q *Q) Write(message []byte) error {
	return q.channel.Publish("", q.queue.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/octet-stream",
			Body:         message,
		},
	)
}
