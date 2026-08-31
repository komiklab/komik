package httphandler

import (
	"net/http"

	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/channels"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
	"github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"
)

// GetChannel implements [apidefn.ServerInterface].
func (h *HttpHandler) GetChannel(ctx *echo.Context) error {
	chans, err := channels.GetChannelList(ctx.Request().Context(), h.dbclient)
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	return ctx.JSON(http.StatusOK, chans)
}

// PostChannel implements [apidefn.ServerInterface].
func (h *HttpHandler) PostChannel(ctx *echo.Context) error {
	var req apidefn.ChannelRequest
	if err := ctx.Bind(&req); utils.IsErrNotNil(err) {
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to bind request", err))
	}
	channelModel, err := models.CreateChannelFromRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create channel model")
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to create channel model", err))
	}
	err = channels.CreateChannel(ctx.Request().Context(), h.dbclient, channelModel)
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	return ctx.JSON(http.StatusOK, channelModel)
}

// PostChannelChannelIdSend implements [apidefn.ServerInterface].
func (h *HttpHandler) PostChannelChannelIdSend(ctx *echo.Context, channelId types.UUID) error {
	channel, err := channels.GetChannelByID(ctx.Request().Context(), h.dbclient, channelId.String())
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	if channel == nil {
		return ctx.JSON(http.StatusNotFound, utils.NewNotFoundError("channel not found", nil))
	}
	//client := hooks.NewSlackWebHookLite(h.cfg.SlackIntegration.BotToken)
	messageSender, err := channels.NewMessageSenderFactory().GetMessageSender(channel)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.NewGeneralError(err))
	}
	err = messageSender.SendMessage(channel, "This is a test message from KomikLab API")
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Message sent successfully"})
}
