package main

import (
	"database/sql"
	"log"
	"webster/events/client"
	"webster/events/handler"
	"webster/events/job"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://webster:gabriel@localhost:5432/websterdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	connection, channel := client.InitRabbitMQ()
	defer connection.Close()
	defer channel.Close()

	r := gin.Default()

	r.POST("/one_time_event", handler.CreateOneTimeEvent(db))
	r.DELETE("/one_time_event/:id", handler.DeleteOneTimeEvent(db))

	r.POST("/recurring_event", handler.CreateRecurringEvent(db))
	r.DELETE("/recurring_event/:id", handler.DeleteRecurringEvent(db))

	go job.StartGenerateEventOccurrencesJob(db)
	go job.SendEventOccurrencesJob(db, channel)

	r.Run(":8080")
}
