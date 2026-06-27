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
