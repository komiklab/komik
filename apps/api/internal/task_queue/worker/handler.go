package worker

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/entity"
	"github.com/komiklab/komik/internal/hooks"
	"github.com/komiklab/komik/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
)

type WorkerHandler struct {
	cfg       *internal.Config
	entitySrv *entity.EntityService
}

func NewWorkerHandler(cfg *internal.Config, dbconn *client.PostgresClient) *WorkerHandler {
	entitySrv := entity.NewEntityService(dbconn)
	return &WorkerHandler{
		cfg:       cfg,
		entitySrv: entitySrv,
	}
}

func (w *WorkerHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	log.Debug().Msg("Processing task")
	// map task into slackevents.AppMentionEvent directly
	var ev slackevents.AppMentionEvent
	if err := json.Unmarshal(task.Payload(), &ev); err != nil {
		log.Error().Msgf("failed to unmarshal slackevents.AppMentionEvent: %v", err)
		return err
	}
	log.Debug().Msgf("slackevents.AppMentionEvent unmarshaled: %v", ev)
	// get the user email id from slack user id
	slackwh := hooks.NewSlackWebHook(w.cfg)
	email, err := slackwh.GetUserEmail(ev.User)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user email")
		return err
	}
	headers := task.Headers()
	acknowledgementId := headers["acknowledgement"]
	log.Debug().Msg("acknowledgement is " + acknowledgementId)
	log.Debug().Msgf("user email: %s", email)
	mentionRe := regexp.MustCompile(`<@[A-Z0-9]+>`)
	cleanText := strings.TrimSpace(mentionRe.ReplaceAllString(ev.Text, ""))
	payloadBytes, err := json.Marshal(map[string]string{"text": cleanText})
	if err != nil{
		log.Error().Err(err).Msg("could not marshal payload")
		return err
	}
	err = w.entitySrv.InitiateEntity(ctx, utils.SourceTypeSlack, acknowledgementId, email, payloadBytes)
	if err != nil {
		log.Error().Err(err).Msg("failed to initiate the entity")
		return err
	}
	return nil
}
