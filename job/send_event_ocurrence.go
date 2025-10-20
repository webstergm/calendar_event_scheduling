package job

import (
	"database/sql"
	"time"
	"webster/events/client"
	"webster/events/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendEventOccurrencesJob(db *sql.DB, rabbitmqChannel *amqp.Channel) {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		sendEventOccurrences(db, rabbitmqChannel)
	}
}

func sendEventOccurrences(db *sql.DB, rabbitmqChannel *amqp.Channel) {
	//query db for occurrences that will happen in the next second either start or end time
	currentTime := time.Now().UTC().Truncate(time.Second)
	occurrences := repository.FindEventOccurencesBetweenTimes(db, currentTime, currentTime.Add(1*time.Second))

	//send those events to a rabbitmq queue
	if len(occurrences) == 0 {
		return
	}

	client.PublishEventOccurences(occurrences, rabbitmqChannel)
	//todo: clean occurrences
}
