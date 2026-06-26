package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/hooks"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack/slackevents"
)

type WorkerHandler struct {
	cfg *internal.Config
}

func NewWorkerHandler(cfg *internal.Config) *WorkerHandler {
	return &WorkerHandler{
		cfg: cfg,
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
	return nil
}
