package entity

const (
	SourceTypeSlack string = "slack"
	SourceTypeAPI   string = "api"
)

func IsValidSourceType(sourceType string) bool {
	switch sourceType {
	case SourceTypeSlack, SourceTypeAPI:
		return true
	default:
		return false
	}
}
