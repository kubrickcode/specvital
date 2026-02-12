package entity

type EventType string

const (
	EventTypeSpecview EventType = "specview"
	EventTypeAnalysis EventType = "analysis"
)

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeSpecview, EventTypeAnalysis:
		return true
	default:
		return false
	}
}
