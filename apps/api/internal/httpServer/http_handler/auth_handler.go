package httphandler

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/auth"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"
	"strings"
	"github.com/rs/zerolog/log"
	// "github.com/coreos/go-oidc/v3/oidc"
	// "net/url"
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

// GetAuthOidcCallback implements [apidefn.ServerInterface].
func (h *HttpHandler) GetAuthOidcCallback(ctx *echo.Context, params apidefn.GetAuthOidcCallbackParams) error {
	authsrv := auth.NewAuthService(h.dbclient, h.cache)
	user, err := authsrv.FetchUser(ctx, params, h.cfg)
	if utils.IsErrNotNil(err) {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
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
	reqId := ctx.Response().Header().Get(echo.HeaderXRequestID)
	msg, err := authsrv.MessageSignInEvent(user.Username, reqId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create sign in event")
	} else {
		err := h.publisher.Publish(auth.EventSigninSubject, msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to publish sign in event")
		}
	}
	return ctx.Redirect(http.StatusFound, h.cfg.PostLoginRedirect)
}

// GetAuthOidcLogin implements [apidefn.ServerInterface].
func (h *HttpHandler) GetAuthOidcLogin(ctx *echo.Context) error {
	oauthConfig := &oauth2.Config{
		ClientID:     h.cfg.OauthConfig.ClientID,
		ClientSecret: h.cfg.OauthConfig.ClientSecret,
		RedirectURL:  h.cfg.OauthConfig.RedirectURL,
		Scopes:       strings.Split(h.cfg.OauthConfig.Scopes, " "),
		Endpoint: oauth2.Endpoint{
			AuthURL:  h.cfg.OauthConfig.AuthURL,
			TokenURL: h.cfg.OauthConfig.TokenURL,
		},
	}
	state := rand.Text()
	// create state cookie
	setCallBackCookie(ctx, "state", state)
	nonce := rand.Text()
	setCallBackCookie(ctx, "nonce", nonce)
	verifier := oauth2.GenerateVerifier()
	setCallBackCookie(ctx, "code_verifier", verifier)
	redirectURL := oauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.S256ChallengeOption(verifier),
	)
	return ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func setCallBackCookie(ctx *echo.Context, name string, value string) error {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   ctx.Request().TLS != nil,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(ctx.Response(), cookie)
	return nil
}
