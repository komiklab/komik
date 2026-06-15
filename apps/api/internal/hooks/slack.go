package hooks

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/komiklab/komik/internal"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type SlackWebHook struct {
	client        *slack.Client
	signingSecret string
}

func NewSlackWebHook(cfg *internal.Config) *SlackWebHook {
	api := slack.New(cfg.SlackIntegration.BotToken, slack.OptionDebug(cfg.IsDebugLoggerConfig))
	return &SlackWebHook{
		client:        api,
		signingSecret: cfg.SlackIntegration.SigningSecret,
	}
}

func (s *SlackWebHook) Handle(ctx *echo.Context, payload []byte) error {
	sv, err := slack.NewSecretsVerifier(ctx.Request().Header, s.signingSecret)
	if err != nil {
		log.Error().Msgf("slack event signature verification setup failed: %v", err)
		return err
	}
	_, err = sv.Write(payload)
	if err != nil {
		log.Error().Msgf("slack event signature verification failed: %v", err)
		return err
	}
	if err := sv.Ensure(); err != nil {
		log.Error().Msgf("slack event signature verification failed: %v", err)
		return err
	}
	eventAPIevent, err := slackevents.ParseEvent(json.RawMessage(payload), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Error().Msgf("failed to parse slack event: %v", err)
		return err
	}
	switch eventAPIevent.Type {
	case slackevents.URLVerification:
		var challenge slackevents.EventsAPIURLVerificationEvent
		if err := json.Unmarshal([]byte(payload), &challenge); err != nil {
			log.Error().Msgf("failed to parse url verification event: %v", err)
			return err
		}
		log.Info().Msgf("url verification event received: %v", challenge.Challenge)
		return ctx.String(http.StatusOK, challenge.Challenge)
	case slackevents.CallbackEvent:
		innerEvent := eventAPIevent.InnerEvent
		mentionRe := regexp.MustCompile(`<@[A-Z0-9]+>`)
		
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			cleanText := strings.TrimSpace(mentionRe.ReplaceAllString(ev.Text, ""))
			log.Info().Msgf("app mention event received: '%s'", cleanText)
		case *slackevents.MessageEvent:
			cleanText := strings.TrimSpace(mentionRe.ReplaceAllString(ev.Text, ""))
			log.Info().Msgf("message event received: '%s'", cleanText)
		default:
			log.Info().Msgf("unknown inner event type: %v", innerEvent.Type)
		}
		return ctx.String(http.StatusOK, "ok")
	default:
		log.Info().Msgf("unknown event type: %v", eventAPIevent.Type)
		return ctx.String(http.StatusOK, "ok")
	}

}
