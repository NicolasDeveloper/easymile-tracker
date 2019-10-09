package broaker

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// IMessagingClient defines our interface for connecting and consuming messages.
type IMessagingClient interface {
	Connect(connectionString string)
	Publish(body []byte, queueName string, routingKey string, exchangeName string, exchangeType string) error
	Close()
}

// AmqpClient is our real implementation, encapsulates a pointer to an amqp.Connection
type AmqpClient struct {
	conn *amqp.Connection
}

// Connect connects to an AMQP broker using the supplied connectionString.
func (m *AmqpClient) Connect(connectionString string) {
	if connectionString == "" {
		panic("Cannot initialize connection to broker, connectionString not set. Have you initialized?")
	}

	var err error
	m.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectionString))
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + connectionString)
	}
}

// Close close connection
func (m *AmqpClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

// Publish publishes a message to the named exchange.
func (m *AmqpClient) Publish(body []byte, queueName string, routingKey string, exchangeName string, exchangeType string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialized. Don't do that.")
	}

	ch, err := m.conn.Channel() // Get a channel from the connection
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	failOnError(err, "Failed to register an Exchange")

	queue, err := ch.QueueDeclare( // Declare a queue that will be created if not exists with some args
		queueName, // our queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	err = ch.QueueBind(
		queue.Name,   // name of the queue
		routingKey,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)

	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body: body,
		})
	handleError(err, string(body))
	return err
}

func failOnError(err error, msg string) {
	if err != nil {
		handleError(err, msg)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
