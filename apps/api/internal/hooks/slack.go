package hooks

import (
	"encoding/json"
	"regexp"
	"net/http"
	"strings"

	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	taskqueue "github.com/komiklab/komik/internal/task_queue"
	asynclient "github.com/komiklab/komik/internal/task_queue/client"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type SlackWebHook struct {
	client        *slack.Client
	signingSecret string
	taskQ         *asynclient.AsynqClient
}

func NewSlackWebHookLite(token string) *SlackWebHook {
	api := slack.New(token, slack.OptionDebug(true))
	return &SlackWebHook{
		client: api,
	}
}

func NewSlackWebHook(cfg *internal.Config) *SlackWebHook {
	redisclient := client.NewRedisClient(cfg)
	taskQ := asynclient.NewAsyncClient(redisclient)
	api := slack.New(cfg.SlackIntegration.BotToken, slack.OptionDebug(cfg.IsDebugLoggerConfig))
	return &SlackWebHook{
		client:        api,
		signingSecret: cfg.SlackIntegration.SigningSecret,
		taskQ:         taskQ,
	}
}

func (s *SlackWebHook) GetUserEmail(userId string) (string, error) {
	user, err := s.client.GetUserInfo(userId)
	if err != nil {
		log.Error().Msgf("failed to get user info: %v", err)
		return "", err
	}
	return user.Profile.Email, nil
}

func (s *SlackWebHook) SendMessage(ev *slackevents.AppMentionEvent, message string) error {
	log.Debug().Msg("channel is " + ev.Channel + "and TimeStamp is " + ev.TimeStamp)
	_, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText(message, false), slack.MsgOptionTS(ev.TimeStamp))
	if err != nil {
		log.Error().Msgf("failed to post message: %v", err)
		return err
	}
	return nil
}

func (s *SlackWebHook) SendMessageLite(channelId string, message string) error {
	_, _, err := s.client.PostMessage(channelId, slack.MsgOptionText(message, false))
	if err != nil {
		log.Error().Msgf("failed to post message: %v", err)
		return err
	}
	return nil
}

func (s *SlackWebHook) Handle(ctx *echo.Context, payload []byte) (*slackevents.AppMentionEvent, error) {
	c := ctx.Request().Context()
	log := log.Ctx(c)
	sv, err := slack.NewSecretsVerifier(ctx.Request().Header, s.signingSecret)
	if err != nil {
		log.Error().Msgf("slack event signature verification setup failed: %v", err)
		return nil, err
	}
	_, err = sv.Write(payload)
	if err != nil {
		log.Error().Msgf("slack event signature verification failed: %v", err)
		return nil, err
	}
	if err := sv.Ensure(); err != nil {
		log.Error().Msgf("slack event signature verification failed: %v", err)
		return nil, err
	}
	eventAPIevent, err := slackevents.ParseEvent(json.RawMessage(payload), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Error().Msgf("failed to parse slack event: %v", err)
		return nil, err
	}
	switch eventAPIevent.Type {
	case slackevents.URLVerification:
		var challenge slackevents.EventsAPIURLVerificationEvent
		if err := json.Unmarshal([]byte(payload), &challenge); err != nil {
			log.Error().Msgf("failed to parse url verification event: %v", err)
			return nil, err
		}
		log.Info().Msgf("url verification event received: %v", challenge.Challenge)
		return nil, ctx.String(http.StatusOK, challenge.Challenge)
	case slackevents.CallbackEvent:
		innerEvent := eventAPIevent.InnerEvent
		mentionRe := regexp.MustCompile(`<@[A-Z0-9]+>`)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			//cleanText := strings.TrimSpace(mentionRe.ReplaceAllString(ev.Text, ""))
			//log.Info().Msgf("app mention event received: '%s'", cleanText)
			// payloadBytes, err := json.Marshal(map[string]string{"text": cleanText})
			// if err != nil {
			// 	log.Error().Msgf("failed to marshal clean text: %v", err)
			// 	return err
			// }
			// generate acknowledgement ID

			ackid := ctx.Response().Header().Get(echo.HeaderXRequestID)
			// add to task queue
			eventBytes, err := json.Marshal(ev)
			//err = fmt.Errorf("some random error")
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal event")
				return ev, err
			}
			headers := map[string]string{
				"acknowledgement": ackid}

			if err := s.taskQ.EnqueueTask(taskqueue.SlackAppMention, eventBytes, headers); err != nil {
				log.Error().Msgf("failed to enqueue task: %v", err)
				return ev, err
			}
			// // send acknowledgement to slack
			// _, _, err = s.client.PostMessage(ev.Channel, slack.MsgOptionText("ACK:"+ackid, false), slack.MsgOptionTS(ev.TimeStamp))
			// if err != nil {
			// 	log.Error().Msgf("failed to post message: %v", err)
			// 	return nil, err
			// }
			//s.entityService.InitiateEntity(entity.SourceTypeSlack, ev.Channel, user.Profile.Email, payloadBytes)
			return ev, ctx.String(http.StatusOK, "ok")
		case *slackevents.MessageEvent:
			cleanText := strings.TrimSpace(mentionRe.ReplaceAllString(ev.Text, ""))
			log.Info().Msgf("message event received: '%s'", cleanText)
		default:
			log.Info().Msgf("unknown inner event type: %v", innerEvent.Type)
		}
		return nil, ctx.String(http.StatusOK, "ok")
	default:
		log.Info().Msgf("unknown event type: %v", eventAPIevent.Type)
		return nil, ctx.String(http.StatusOK, "ok")
	}

}
