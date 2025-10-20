package eventtype

type EventType string

const (
	OneTime   EventType = "one_time"
	Recurring EventType = "recurring"
)
