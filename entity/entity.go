package entity

import "time"

type RecurrenceRule struct {
	Id             int64
	EventId        int64
	StartDate      time.Time
	EndDate        time.Time
	RRule          string
	RRuleStartTime time.Time
	RRuleEndTime   time.Time
}
