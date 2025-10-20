package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"webster/events/model"
	eventtype "webster/events/model/enum"
	"webster/events/repository"

	"github.com/gin-gonic/gin"
)

func CreateOneTimeEvent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload model.OneTimeEvent
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		eventID, err := createOneTimeEventWithTransaction(db, &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "One-time event created successfully",
			"event_id": eventID,
		})
	}
}

func createOneTimeEventWithTransaction(db *sql.DB, payload *model.OneTimeEvent) (int64, error) {
	var err error
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	var eventId int64
	eventId, err = repository.InsertEventWithTx(tx, payload.Name, eventtype.OneTime)
	if err != nil {
		return 0, err
	}

	start := payload.StartDate.UTC()
	end := payload.EndDate.UTC()

	if err = repository.InsertEventOccurrenceWithTx(tx, &model.EventOccurence{EventId: eventId, StartDate: start, EndDate: end}); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return eventId, nil
}

func DeleteOneTimeEvent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM events WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "One-time event deleted"})
	}
}
