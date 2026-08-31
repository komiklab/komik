package channels

import (
	"encoding/json"
	"fmt"

	"github.com/komiklab/komik/internal/hooks"
	"github.com/komiklab/komik/internal/models"
)

// A factory function to send message. As of now it only supports Slack
// note config will be different for each channel

type MessageSender interface {
	SendMessage(channel *models.Channel, message string) error
}

type MessageSenderFactory struct {
}

func NewMessageSenderFactory() *MessageSenderFactory {
	return &MessageSenderFactory{}
}

func (f *MessageSenderFactory) GetMessageSender(channel *models.Channel) (MessageSender, error) {
	switch channel.Type {
	case models.ChannelTypeSlack:
		return NewSlackMessageSender(channel)
	default:
		return nil, fmt.Errorf("unsupported channel type: %s", channel.Type)
	}
}

type SlackMessageSender struct {
	hook *hooks.SlackWebHook
	channelID string

}

func NewSlackMessageSender(channel *models.Channel) (*SlackMessageSender, error) {
	var slackConfigStruct struct {
		BotToken string 
		ChannelName string 
	}
	slackConfig := channel.Config
	err := json.Unmarshal(slackConfig, &slackConfigStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Slack config: %w", err)
	}
	hook := hooks.NewSlackWebHookLite(slackConfigStruct.BotToken)
	return &SlackMessageSender{
		hook: hook,
		channelID: slackConfigStruct.ChannelName,
	}, nil
}

func (s *SlackMessageSender) SendMessage(channel *models.Channel, message string) error {
	// Implement the logic to send a message to Slack using the channel's config
	// You can use the Slack API client here
	return s.hook.SendMessageLite(s.channelID, message)
}

