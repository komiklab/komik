package httphandler

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
)

type SessionAuthenticator struct {
	redisclient *client.RedisClient
}

func NewSessionAuthenticator(redisclient *client.RedisClient) *SessionAuthenticator {
	return &SessionAuthenticator{
		redisclient: redisclient,
	}
}

func (s *SessionAuthenticator) Authenticate(ctx *echo.Context) error {
	// 1. Check for session cookie
	cookie, err := ctx.Cookie(utils.SESSION_COOKIE_NAME)
	if err != nil {
		return utils.NewAuthenticationError("missing session cookie", err)
	}
	
	// 2. Validate with Redis
	sessionID := cookie.Value
	sessionData, err := s.redisclient.Get("session:" + sessionID)
	if err != nil {
		return utils.NewAuthenticationError("invalid or expired session", err)
	}
	
	// 3. Inject the session data into the context
	ctx.Set("user_session", sessionData)
	ctx.Set("session_id", "session:"+sessionID)
	
	return nil
}