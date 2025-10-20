package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	"webster/events/model"
)

func InsertEventOccurrenceWithTx(tx *sql.Tx, eventOccurence *model.EventOccurence) error {
	query := `
		INSERT INTO event_occurrences(event_id, datetime_start, datetime_end)
		VALUES($1, $2, $3)
	`

	_, err := tx.Exec(query, eventOccurence.EventId, eventOccurence.StartDate, eventOccurence.EndDate)
	if err != nil {
		return fmt.Errorf("failed to insert event occurrence: %w", err)
	}

	return nil
}

func InsertEventOccurrencesWithTx(tx *sql.Tx, eventOccurences []model.EventOccurence) error {
	query := `
		INSERT INTO event_occurrences(event_id, datetime_start, datetime_end) VALUES 
	`
	args := []interface{}{}
	placeholders := []string{}

	for i, eventOccurence := range eventOccurences {
		n := i*3 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", n, n+1, n+2))
		args = append(args, eventOccurence.EventId, eventOccurence.StartDate, eventOccurence.EndDate)
	}

	query += strings.Join(placeholders, ",")

	_, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert event occurrences: %w", err)
	}

	return nil
}

func FindEventOccurencesBetweenTimes(db *sql.DB, startTime time.Time, endTime time.Time) []model.EventOccurence {
	var eventOccurrences []model.EventOccurence
	query := `
		SELECT event_id, datetime_start, datetime_end FROM event_occurrences 
			WHERE (datetime_start >= $1 AND datetime_start < $2)
			OR (datetime_end >= $1 AND datetime_end < $2)
	`
	rows, err := db.Query(query, startTime, endTime)
	// extra logging
	log.Printf("Args: %v, %v", startTime, endTime)
	if err != nil {
		log.Printf("Error while querying event occurrences: %v", err)
		return eventOccurrences
	}
	defer rows.Close()

	for rows.Next() {
		var e model.EventOccurence
		if err := rows.Scan(
			&e.EventId,
			&e.StartDate,
			&e.EndDate,
		); err != nil {
			log.Printf("commit failed: %v", err)
			return eventOccurrences
		}
		eventOccurrences = append(eventOccurrences, e)
	}

	// extra logging
	if len(eventOccurrences) > 0 {
		log.Printf("Event occurrences found: %d", len(eventOccurrences))
		for i, eventOccurrence := range eventOccurrences {
			log.Printf("event_occurrence[%d]: event_id = %d, start = %v, end = %v",
				i, eventOccurrence.EventId, eventOccurrence.StartDate, eventOccurrence.EndDate)
		}
	}

	return eventOccurrences
}
