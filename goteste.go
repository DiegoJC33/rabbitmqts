package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

func main() {
	// Connection information
	conf := map[string]string{
		"host":       "dns",
		"port":       "5671",
		"queue_name": "user",
		"username":   "queue",
		"password":   "pass",
	}

	// Configure the TLS connection
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	// Connect to the RabbitMQ server with TLS
	conn, err := amqp.DialTLS("amqps://"+conf["username"]+":"+conf["password"]+"@"+conf["host"]+":"+conf["port"]+"/", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(conf["queue_name"], true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Send a message to the queue
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello, World!"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message sent")

	// Capture interrupt signal (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-interrupt

	fmt.Println("Closing connection...")
}
