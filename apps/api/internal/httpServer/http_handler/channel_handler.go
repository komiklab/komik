package httphandler

import (
	"net/http"

	"github.com/komiklab/komik/internal/channels"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
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
	panic("unimplemented")
}
