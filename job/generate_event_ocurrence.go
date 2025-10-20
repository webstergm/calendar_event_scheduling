package job

import (
	"database/sql"
	"log"
	"time"
	"webster/events/entity"
	"webster/events/model"
	"webster/events/repository"

	"github.com/teambition/rrule-go"
)

var pageSize uint32 = 30

func StartGenerateEventOccurrencesJob(db *sql.DB) {
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		generateEventOccurrences(db)
	}
}

func generateEventOccurrences(db *sql.DB) {
	currentTime := time.Now()
	pageNr := uint32(0)

	for {
		tx, err := db.Begin()
		if err != nil {
			log.Printf("failed to begin tx: %v", err)
			return
		}

		defer tx.Rollback()

		rules, err := repository.SelectRecurrenceRulesForOcurrenceGenerationWithTx(tx, currentTime, pageSize, pageNr)
		if err != nil {
			log.Printf("select recurrence rules failed: %v", err)
			return
		}

		for _, r := range rules {
			occurrences := generateOccurrencesForRule(r)
			if err = repository.InsertEventOccurrencesWithTx(tx, occurrences); err != nil {
				log.Printf("insert occurrences failed: %v", err)
				return
			}
		}

		if err = tx.Commit(); err != nil {
			log.Printf("commit failed: %v", err)
			return
		}

		if len(rules) <= int(pageSize) {
			break
		}
		pageNr++
	}
}

func generateOccurrencesForRule(r entity.RecurrenceRule) []model.EventOccurence {
	from := time.Now()
	until := from.AddDate(0, 0, 10)
	if until.After(r.EndDate) {
		until = r.EndDate.AddDate(0, 0, 1)
	}

	rule, _ := rrule.StrToRRule(r.RRule)
	times := rule.Between(from, until, true)

	occurrences := make([]model.EventOccurence, len(times))
	for i, t := range times {
		occurrences[i] = model.EventOccurence{
			EventId:   r.EventId,
			StartDate: time.Date(t.Year(), t.Month(), t.Day(), r.RRuleStartTime.Hour(), r.RRuleStartTime.Minute(), r.RRuleStartTime.Second(), 0, t.Location()),
			EndDate:   time.Date(t.Year(), t.Month(), t.Day(), r.RRuleEndTime.Hour(), r.RRuleEndTime.Minute(), r.RRuleEndTime.Second(), 0, t.Location()),
		}
	}
	return occurrences
}
