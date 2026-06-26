package taskqueue

import (
	"errors"
)

type TaskType string

const (
	SlackAppMention TaskType = "slack:appmention"
)

func ParseTaskType(t TaskType) (TaskType, error) {
	switch t {
	case SlackAppMention:
		return SlackAppMention, nil
	default:
		return "", errors.New("invalid task type: " + string(t))
	}
}