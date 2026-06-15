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
	//var req apidefn.PostHookSlackJSONRequestBody
	// first convert the request body in []bytes
	bodyBytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Msgf("error reading request body: %v", err)
		return err
	}
	// create a slack web hook
	slackHook := hooks.NewSlackWebHook(h.cfg)
	// handle slack event
	if err := slackHook.Handle(ctx, bodyBytes); err != nil {
		log.Error().Msgf("error handling slack event: %v", err)
		return err
	}
	return nil
}
