package repository

import (
	"database/sql"
	"fmt"
	"time"
	"webster/events/entity"
)

func InsertRecurrenceRuleWithTx(tx *sql.Tx, eventId int64, start time.Time, end time.Time,
	rrule string, rruleStartTime string, rruleEndTime string) error {

	_, err := tx.Exec(`
        INSERT INTO recurrence_rules(event_id, start_date, end_date, rrule, rrule_start_time, rrule_end_time)
        VALUES($1, $2, $3, $4, $5, $6)
    `, eventId, start, end, rrule, rruleStartTime, rruleEndTime)
	if err != nil {
		return fmt.Errorf("failed to insert recurrence rule: %w", err)
	}

	return nil
}

func SelectRecurrenceRulesForOcurrenceGenerationWithTx(tx *sql.Tx, currentTime time.Time, pageSize uint32, pageCount uint32) ([]entity.RecurrenceRule, error) {
	rows, err := tx.Query(`
		SELECT id, event_id, start_date, end_date, rrule, rrule_start_time, rrule_end_time FROM recurrence_rules rr
			WHERE $1::timestamp BETWEEN start_date AND end_date
			AND NOT EXISTS (
    			SELECT 1
    			FROM event_occurrences eo
    		WHERE rr.event_id = eo.event_id
		)
		ORDER BY rr.id
		LIMIT $2 OFFSET $3;
	`, currentTime, pageSize+1, pageCount*pageSize)

	var recurrenceRules []entity.RecurrenceRule
	if err != nil {
		return recurrenceRules, err
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.RecurrenceRule
		if err := rows.Scan(
			&r.Id,
			&r.EventId,
			&r.StartDate,
			&r.EndDate,
			&r.RRule,
			&r.RRuleStartTime,
			&r.RRuleEndTime,
		); err != nil {
			return nil, err
		}
		recurrenceRules = append(recurrenceRules, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recurrenceRules, nil
}
