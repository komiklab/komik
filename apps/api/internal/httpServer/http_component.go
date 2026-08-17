package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator/v12"
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/event"
	httphandler "github.com/komiklab/komik/internal/httpServer/http_handler"
	"github.com/komiklab/komik/internal/utils"

	// httphandler "github.com/komiklab/komik/internal/httpServer/handler/http_handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog/log"
)

type HttpComponent struct {
	e           *echo.Echo
	s           *http.Server
	cfg         *internal.Config
	dbclient    *client.PostgresClient
	publisher   *event.Publisher
	redisclient *client.RedisClient
}

var _ internal.Component = (*HttpComponent)(nil)

func (h *HttpComponent) GetName() string {
	return "HttpComponent"
}

func (h *HttpComponent) Init() {
	h.e = echo.New()
	h.e.Validator = &CustomValidator{}
	h.addMiddleware()

	sessionAuth := httphandler.NewSessionAuthenticator(h.redisclient)
	hmacAuth := httphandler.NewHmacAuthenticator(h.dbclient)

	handler := httphandler.NewHttpHandler(h.cfg, h.dbclient, h.publisher, h.redisclient)
	apidefn.RegisterHandlersWithOptions(h.e, handler, apidefn.RegisterHandlersOptions{
		BaseURL: "/api/v1",
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"GetAuthMe":      {h.AuthMiddleware(sessionAuth)},
			"PostAuthLogout": {h.AuthMiddleware(sessionAuth)},
			"GetHook":        {h.AuthMiddleware(sessionAuth)},
			"PostHook":       {h.AuthMiddleware(sessionAuth)},
			"PostHookId":     {h.AuthMiddleware(hmacAuth)},
		},
	})
}

func (h *HttpComponent) Start() {
	h.s = &http.Server{
		Addr:              ":65080",
		Handler:           h.e,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	if err := h.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("HTTP server failed to start")
	}
}

func (h *HttpComponent) Stop(ctx context.Context) {
	defer h.dbclient.Close()
	if h.s != nil {
		if err := h.s.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("HTTP server failed to shutdown")
		}
	}
}

func NewHttpComponent(cfg *internal.Config, dbclient *client.PostgresClient, publisher *event.Publisher, redisclient *client.RedisClient) *HttpComponent {
	return &HttpComponent{
		cfg:         cfg,
		dbclient:    dbclient,
		publisher:   publisher,
		redisclient: redisclient,
	}
}

func (h *HttpComponent) addMiddleware() {
	h.e.Use(middleware.Recover())
	h.e.Use(middleware.RequestID())
	h.e.Use(h.LoggerContextMiddleware())
	h.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{h.cfg.CORSSupport},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			// required so CORS preflight allows the CSRF header
			"X-CSRF-Token",
			"X-Signature",
			"X-Timestamp", 
		},
		// ExposeHeaders allows the browser to expose these response headers to
		// frontend JavaScript on cross-origin requests.
		ExposeHeaders:    []string{"X-CSRF-Token"},
		AllowCredentials: true,
	}))
	h.e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper: func(c *echo.Context) bool {
			// skip the csrf for /webhook paths
			return c.Path() == "/api/v1/hook/slack" //|| c.Path() == "/api/v1/hook/:id"
		},
		TokenLookup:    "header:X-CSRF-Token",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		TrustedOrigins: []string{h.cfg.CORSSupport},
		CookieSecure:   false, // must be false for local HTTP dev; set to true behind TLS in prod
		CookieHTTPOnly: false, // must be false so JS can read the cookie value
		CookieSameSite: http.SameSiteLaxMode,
	}))
	h.e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogMethod:    true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("request_id", v.RequestID).
				Str("method", v.Method).
				Msg("request")

			return nil
		},
	}))
}

type CustomValidator struct {
}

func (c *CustomValidator) Validate(i any) error {
	_, err := govalidator.ValidateStruct(i)
	if err != nil {
		return utils.NewValidationError("validation failed", err)
	}
	return nil
}

func (h *HttpComponent) AuthMiddleware(authenticators httphandler.Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			var lastErr error
			err := authenticators.Authenticate(ctx)
			if err == nil {
				return next(ctx)
			}
			lastErr = err
			return ctx.JSON(http.StatusUnauthorized, lastErr)
		}
	}
}

func (h *HttpComponent) LoggerContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = c.Request().Header.Get(echo.HeaderXRequestID)
			}

			// Create a child logger with the request_id
			logger := log.With().Str("request_id", reqID).Logger()

			// Inject logger into the request context
			req := c.Request()
			ctx := logger.WithContext(req.Context())
			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}
