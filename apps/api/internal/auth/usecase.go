package auth

import (
	"errors"
	"net/url"

	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	// "golang.org/x/oauth2"
)

type AuthService struct {
	authrepo *repositories.AdminRepo
}

func (a *AuthService) Logout(sessionId string) error {
	err := a.authrepo.DeleteSession(sessionId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete session")
		return err
	}
	return nil
}

func NewAuthService(dbcon *client.PostgresClient, cache *client.RedisClient) *AuthService {
	authrepo := repositories.NewAdminRepo(dbcon, cache)
	return &AuthService{
		authrepo: authrepo,
	}
}

func (a *AuthService) DoesAdminExist() (bool, error) {
	exist, err := a.authrepo.DoesAdminExist()
	if err != nil {
		return false, utils.NewGeneralError(err)
	}
	return exist, nil
}

func (a *AuthService) CreateAdmin(admin *models.Admin) error {
	var err error
	admin.Password, err = utils.HashPassword(admin.Password)
	if err != nil {
		return utils.NewGeneralError(err)
	}
	return a.authrepo.CreateAdmin(admin)
}

func (a *AuthService) VerifyPassword(admin *models.Admin) error {
	passwordReceived := admin.Password
	hashedPassword, err := a.authrepo.FetchPassword(admin)
	if utils.IsErrNotNil(err) {
		return err
	}
	verified, err := utils.VerifyPassword(passwordReceived, hashedPassword)
	if utils.IsErrNotNil(err) {
		return utils.NewInternalServerError("could not verify password because of internal issue", err)
	}
	if !verified {
		return utils.NewAuthenticationError("invalid password", errors.New("wrong password"))
	}
	return nil
}

func (a *AuthService) CreateSession(user models.UserRepresentation) (string, error) {
	// we will first store the user if its new in database
	err := a.authrepo.SaveUserIfNotExist(user)
	if utils.IsErrNotNil(err) {
		return "", err
	}
	// then we will create session for the user
	sessionID, err := a.authrepo.CreateSession(&user)
	if utils.IsErrNotNil(err) {
		return "", err
	}
	return sessionID, nil
}

func (a *AuthService) FetchUser(ctx *echo.Context, params apidefn.GetAuthOidcCallbackParams, cfg *internal.Config) (*models.UserRepresentation, error) {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.OauthConfig.ClientID,
		ClientSecret: cfg.OauthConfig.ClientSecret,
		RedirectURL:  cfg.OauthConfig.RedirectURL,
		Scopes:       strings.Split(cfg.OauthConfig.Scopes, " "),
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.OauthConfig.AuthURL,
			TokenURL: cfg.OauthConfig.TokenURL,
		},
	}
	state, err := ctx.Cookie("state")
	if utils.IsErrNotNil(err) {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}
	queryState := ctx.QueryParam("state")
	if utils.IsErrNotNil(err) {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}
	if queryState != state.Value {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}
	codeVerifier, err := ctx.Cookie("code_verifier")
	if utils.IsErrNotNil(err) {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}
	nonce, err := ctx.Cookie("nonce")
	if utils.IsErrNotNil(err) {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}

	token, err := oauthConfig.Exchange(ctx.Request().Context(), params.Code, oauth2.VerifierOption(codeVerifier.Value))
	if utils.IsErrNotNil(err) {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to extract session from context"))
	}
	// 	// verify ID token
	authURL, err := url.Parse(cfg.OauthConfig.AuthURL)
	if err != nil {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to parse auth url"))
	}
	issuerURL := authURL.Scheme + "://" + authURL.Host + "/"

	provider, err := oidc.NewProvider(ctx.Request().Context(), issuerURL)
	if err != nil {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("unable to create oidc provider"))
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.OauthConfig.ClientID})
	idToken, err := verifier.Verify(ctx.Request().Context(), rawIDToken)
	if err != nil {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("invalid id token"))
	}
	// check id token nonce
	if idToken.Nonce != nonce.Value {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("invalid nonce"))
	}
	// check id token iss
	if idToken.Issuer != issuerURL {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("invalid issuer"))
	}
	// check id token sub
	if idToken.Subject == "" {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("invalid subject"))
	}

	// get the user now
	userInfo, err := provider.UserInfo(ctx.Request().Context(), oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, utils.NewAuthenticationError("unauthorized", errors.New("failed to get user info"))
	}
	user := models.NewUserRepresentation(userInfo.Email)
	return user, nil

}

func (a *AuthService) RetrieveUserFromSession(sessionId string) (*models.UserRepresentation, error) {
	return a.authrepo.FetchUserFromSession(sessionId)
}
