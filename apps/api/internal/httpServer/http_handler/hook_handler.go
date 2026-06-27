package httphandler

import (
	//"net/http"

	"io"

	//"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/hooks"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

// PostHookSlack implements [apidefn.ServerInterface].
func (h *HttpHandler) PostHookSlack(ctx *echo.Context) error {
	newctx := ctx.Request().Context()
	log := log.Ctx(newctx)
	// get the request id from ctx
	requestID := ctx.Response().Header().Get(echo.HeaderXRequestID)
	log.Debug().Msg("request id : "+requestID)
	// first convert the request body in []bytes
	bodyBytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Msgf("error reading request body: %v", err)
		return err
	}
	// create a slack web hook
	slackHook := hooks.NewSlackWebHook(h.cfg)
	// handle slack event
	ev, err := slackHook.Handle(ctx, bodyBytes); 
	if err != nil {
		log.Error().Msgf("error handling slack event: %v", err)
		// at this time we will send in slack an error message
		slacksendingerr := slackHook.SendMessage(ev, "error handling slack event : "+requestID)
		if slacksendingerr != nil {
			log.Error().Msgf("error sending slack message: %v", slacksendingerr)
		}
		return err
	}
	// at this time we will send a ack
	slacksendingerr := slackHook.SendMessage(ev, "acknowledged with ID : "+requestID)
	if slacksendingerr != nil {
		log.Error().Msgf("error sending slack message: %v", slacksendingerr)
	}
	return nil
}
