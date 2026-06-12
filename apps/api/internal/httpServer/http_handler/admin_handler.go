package httphandler

import (
	"net/http"

	"github.com/asaskevich/govalidator/v12"
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/auth"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
)

// GetAdmin implements [apidefn.ServerInterface].
func (h *HttpHandler) GetAdmin(ctx *echo.Context) error {

	// if csrfToken, ok := (*ctx).Get("csrf").(string); ok {
	// 	(*ctx).Response().Header().Set("X-CSRF-Token", csrfToken)
	// }

	auth := auth.NewAuthService(h.dbclient)
	exists, err := auth.DoesAdminExist()
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	if exists {
		return ctx.JSON(200, map[string]interface{}{
			"exists": true,
		})
	}
	return ctx.JSON(200, map[string]interface{}{
		"exists": false,
	})
}

// PostAdmin implements [apidefn.ServerInterface].
func (h *HttpHandler) PostAdmin(ctx *echo.Context) error {
	var req apidefn.AdminCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to bind request", err))
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.NewValidationError("validation failed", err))
	}
	auth := auth.NewAuthService(h.dbclient)
	adminDao := models.NewAdminDAO(req.Username, req.Password)
	_, err := govalidator.ValidateStruct(adminDao)
	if utils.IsErrNil(err) {
		return ctx.JSON(http.StatusBadRequest, utils.NewValidationError("validation failed", err))
	}
	if err := auth.CreateAdmin(models.NewAdminDAO(req.Username, req.Password)); err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	return ctx.NoContent(http.StatusCreated)
}
