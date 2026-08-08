package httphandler

import (
	//"net/http"

	"errors"
	"io"
	"net/http"

	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/hooks"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

// GetHook implements [apidefn.ServerInterface].
func (h *HttpHandler) GetHook(ctx *echo.Context) error {
	svc := hooks.NewHookService(h.dbclient)
	hooks, err := svc.FetchHooks()
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}

	return ctx.JSON(http.StatusOK, &models.ListHooks{Hooks: hooks})
}

// PostHook implements [apidefn.ServerInterface].
func (h *HttpHandler) PostHook(ctx *echo.Context) error {
	var req apidefn.HookRegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to bind request", err))
	}
	if req.Name == "" {
		return ctx.JSON(http.StatusBadRequest, utils.NewValidationError("name is missing", errors.New("name is missing")))
	}
	hook := models.NewHook(req.Name)
	svc := hooks.NewHookService(h.dbclient)
	err := svc.CreateHook(hook)
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	return ctx.JSON(http.StatusOK, hook)
}

// PostHookSlack implements [apidefn.ServerInterface].
func (h *HttpHandler) PostHookSlack(ctx *echo.Context) error {
	newctx := ctx.Request().Context()
	log := log.Ctx(newctx)
	// get the request id from ctx
	requestID := ctx.Response().Header().Get(echo.HeaderXRequestID)
	log.Debug().Msg("request id : " + requestID)
	// first convert the request body in []bytes
	bodyBytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Msgf("error reading request body: %v", err)
		return err
	}
	// create a slack web hook
	slackHook := hooks.NewSlackWebHook(h.cfg)
	// handle slack event
	ev, err := slackHook.Handle(ctx, bodyBytes)
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
