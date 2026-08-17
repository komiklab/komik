package models

import (
	"encoding/json"

	"github.com/komiklab/komik/apidefn"
)

type HookInput struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Message   []byte `json:"message"`
	Initiator string `json:"initiator"`
}

func NewHookInputFromPayload(hookname string, req *apidefn.HookSendRequest) (*HookInput, error) {
	initiator := ""
	if req.Initiator != nil {
		initiator = *req.Initiator
	}
	// convert req.Message to []bytes
	msg, err := json.Marshal(req.Message)
	if err != nil {
		return nil, err
	}
	return &HookInput{
		Name:      hookname,
		Reference: req.Reference,
		Message:   msg,
		Initiator: initiator,
	}, nil
}
