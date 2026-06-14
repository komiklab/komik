package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
)

// GetAuthMe implements [apidefn.ServerInterface].
func (h *HttpHandler) GetAuthMe(ctx *echo.Context) error {
	userData := ctx.Get("user_session")
	if userData == nil {
		return utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}

	// Redis returns a JSON string, so userData is a string
	sessionStr, ok := userData.(string)
	if !ok {
		return utils.NewInternalServerError("unable to extract session from context", errors.New("invalid session data format"))
	}

	// Unmarshal the JSON string back into the UserRepresentation struct
	var userRep models.UserRepresentation
	if err := json.Unmarshal([]byte(sessionStr), &userRep); err != nil {
		return utils.NewInternalServerError("failed to parse session data", err)
	}

	return ctx.JSON(http.StatusOK, userRep)
}
