package httphandler

import (
	"net/http"

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
		return ctx.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
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
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError(err.Error()))
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.NewValidationError(err.Error()))
	}
	auth := auth.NewAuthService(h.dbclient)
	if err := auth.CreateAdmin(&models.Admin{Username: req.Username, Password: req.Password}); err != nil {
		return ctx.JSON(err.StatusCode, utils.NewInternalServerError(err.Error()))
	}
	return ctx.NoContent(http.StatusCreated)
}
