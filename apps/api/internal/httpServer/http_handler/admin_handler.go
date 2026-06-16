package httphandler

import (
	"errors"
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

	auth := auth.NewAuthService(h.dbclient, h.cache)
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
	if err := ctx.Bind(&req); utils.IsErrNotNil(err) {
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to bind request", err))
	}
	authsrv := auth.NewAuthService(h.dbclient, h.cache)
	adminDao := models.NewAdminDAO(req.Username, req.Password)
	_, err := govalidator.ValidateStruct(adminDao)
	if utils.IsErrNotNil(err) {
		return ctx.JSON(http.StatusBadRequest, utils.NewValidationError("validation failed", err))
	}
	if err := authsrv.CreateAdmin(models.NewAdminDAO(req.Username, req.Password)); utils.IsErrNotNil(err) {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	// publish an successful event
	reqId := ctx.Response().Header().Get(echo.HeaderXRequestID)
	msg, err := authsrv.MessageEventAdminCreated(adminDao.Username, reqId)
	if !utils.IsErrNotNil(err) {
		err := h.publisher.Publish(auth.EventAdminCreatedSubject, msg)
		utils.IsErrNotNil(err)
	}
	return ctx.NoContent(http.StatusCreated)
}

// PostAuthLogin implements [apidefn.ServerInterface].
func (h *HttpHandler) PostAuthLogin(ctx *echo.Context) error {
	var req apidefn.AdminCreateRequest
	if err := ctx.Bind(&req); utils.IsErrNotNil(err) {
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to bind request", err))
	}
	authsrv := auth.NewAuthService(h.dbclient, h.cache)
	adminDao := models.NewAdminDAO(req.Username, req.Password)
	if err := authsrv.VerifyPassword(adminDao); utils.IsErrNotNil(err) {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	user := models.NewUserRepresentation(adminDao.Username)
	sessionID, err := authsrv.CreateSession(*user)
	if utils.IsErrNotNil(err) {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	// lets set http cookie
	sessionTTL := utils.SESSION_TTL
	sessionCookie := http.Cookie{
		Name:     utils.SESSION_COOKIE_NAME,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   ctx.Request().TLS != nil,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(ctx.Response(), &sessionCookie)
	//return ctx.Redirect(http.StatusFound, h.cfg.PostLoginRedirect)
	return ctx.NoContent(http.StatusOK)
}

// PostAuthLogout implements [apidefn.ServerInterface].
func (h *HttpHandler) PostAuthLogout(ctx *echo.Context) error {
	sessionId := ctx.Get("session_id").(string)
	authsrv := auth.NewAuthService(h.dbclient, h.cache)
	err := authsrv.Logout(sessionId)
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		var redisNilErr *utils.KomikError
		if errors.As(err, &redisNilErr) {
			komikErr.WithStatusCode(http.StatusPreconditionFailed)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	// invalidate the session cookie
	c := &http.Cookie{
		Name:     utils.SESSION_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   ctx.Request().TLS != nil,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(ctx.Response(), c)
	return ctx.NoContent(http.StatusOK)
}
