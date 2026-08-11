package httphandler

import "github.com/labstack/echo/v5"

type Authenticator interface {
	Authenticate(ctx *echo.Context) error
}
