package handler

import (
	"database/sql"
	"net/http"
	"webster/events/model"
	eventtype "webster/events/model/enum"
	"webster/events/repository"
	"webster/events/util"

	"github.com/gin-gonic/gin"
)

func CreateRecurringEvent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload model.RecurringEvent
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		start := payload.StartDate.UTC()
		end := payload.EndDate.UTC()

		var rruleStart string
		rruleStart, err := util.NormalizeTimeToUTC(payload.RRuleStartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rrule_start_time"})
			return
		}

		var rruleEnd string
		rruleEnd, err = util.NormalizeTimeToUTC(payload.RRuleEndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rrule_end_time"})
			return
		}

		var tx *sql.Tx
		tx, err = db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
			return
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
		eventId, err = repository.InsertEventWithTx(tx, payload.Name, eventtype.Recurring)
		if err != nil {
			return
		}

		err = repository.InsertRecurrenceRuleWithTx(tx, eventId, start, end, payload.RRule, rruleStart, rruleEnd)
		if err != nil {
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Recurring event created successfully",
			"event_id": eventId,
		})
	}
}

func DeleteRecurringEvent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM events WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Recurring event deleted"})
	}
}
