package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCreateRecurringEvent_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// begin transaction
	mock.ExpectBegin()

	// expect event insert
	mock.ExpectQuery(`INSERT INTO events`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// expect recurrence rule insert
	mock.ExpectExec(`INSERT INTO recurrence_rules`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// commit transaction
	mock.ExpectCommit()

	router := gin.Default()
	router.POST("/recurring_event", CreateRecurringEvent(db))

	body := `{
		"name": "Weekly Meeting",
		"start_date": "` + time.Now().Format(time.RFC3339) + `",
		"end_date": "` + time.Now().AddDate(0, 1, 0).Format(time.RFC3339) + `",
		"rrule": "FREQ=WEEKLY;BYDAY=MO",
		"rrule_start_time": "09:00:00",
		"rrule_end_time": "10:00:00"
	}`

	req := httptest.NewRequest(http.MethodPost, "/recurring_event", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Recurring event created successfully")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateRecurringEvent_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, _ := sqlmock.New()
	defer db.Close()

	router := gin.Default()
	router.POST("/recurring_event", CreateRecurringEvent(db))

	// missing closing brace should give invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/recurring_event", strings.NewReader(`{"name": "bad json"`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
