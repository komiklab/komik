package utils

const (
	SourceTypeSlack = "SLACK"
)

func IsValidSourceType(sourceType string) bool {
	switch sourceType {
	case SourceTypeSlack:
		return true
	default:
		return false
	}
}

const (
	SystemInitiator = "SYSTEM"
)

const (
	InitiatedTopic  = "komik.entity.initiated"
	DispatchedTopic = "komik.entity.dispatched"
)
