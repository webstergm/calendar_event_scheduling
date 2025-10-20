package client

import (
	"encoding/json"
	"log"
	"webster/events/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://webster:gabriel@localhost:5672/") // hardcoded but who cares for learning project
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	log.Println("Connected to RabbitMQ")

	return conn, ch
}

func PublishEventOccurences(eventOccurrences []model.EventOccurence, rabbitmqChannel *amqp.Channel) {
	for _, occurrence := range eventOccurrences {
		body, _ := json.Marshal(occurrence)

		err := rabbitmqChannel.Publish(
			"calendar_exchange",
			"calendar_events",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
		}
	}
}
