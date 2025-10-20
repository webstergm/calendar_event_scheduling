package repository

import (
	"database/sql"
	"fmt"
	eventtype "webster/events/model/enum"
)

func InsertEventWithTx(tx *sql.Tx, name string, eventType eventtype.EventType) (int64, error) {
	var eventID int64
	query := `INSERT INTO events(name, type) VALUES($1, $2) RETURNING id`

	err := tx.QueryRow(query, name, eventType).Scan(&eventID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert event: %w", err)
	}

	return eventID, nil
}
