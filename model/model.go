package model

import "time"

type OneTimeEvent struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type RecurringEvent struct {
	Name           string    `json:"name" binding:"required"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`
	RRule          string    `json:"rrule" binding:"required"`
	RRuleStartTime string    `json:"rrule_start_time"`
	RRuleEndTime   string    `json:"rrule_end_time"`
}

type EventOccurence struct {
	EventId   int64
	StartDate time.Time
	EndDate   time.Time
}
