package handler

import (
	"testing"
	"time"
	"webster/events/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func Test_createOneTimeEventWithTransaction_Success(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO events`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO event_occurrences`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	payload := &model.OneTimeEvent{
		Name:      "Test Event",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour),
	}

	//call to tested function
	id, err := createOneTimeEventWithTransaction(db, payload)
	require.NoError(t, err)
	require.Equal(t, int64(1), id)

	require.NoError(t, mock.ExpectationsWereMet())
}
