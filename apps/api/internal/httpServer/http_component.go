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
	e         *echo.Echo
	s         *http.Server
	cfg       *internal.Config
	dbclient  *client.PostgresClient
	publisher *event.Publisher
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
	handler := httphandler.NewHttpHandler(h.cfg, h.dbclient, h.publisher, h.redisclient)
	apidefn.RegisterHandlersWithOptions(h.e, handler, apidefn.RegisterHandlersOptions{
		BaseURL:              "/api/v1",
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"GetAuthMe":      {h.AuthMiddleware()},
			"PostAuthLogout": {h.AuthMiddleware()},
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
		cfg:       cfg,
		dbclient:  dbclient,
		publisher: publisher,
		redisclient: redisclient,
	}
}

func (h *HttpComponent) addMiddleware() {
	h.e.Use(middleware.Recover())
	h.e.Use(middleware.RequestID())
	h.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{h.cfg.CORSSupport},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-CSRF-Token", // required so CORS preflight allows the CSRF header
		},
		// ExposeHeaders allows the browser to expose these response headers to
		// frontend JavaScript on cross-origin requests.
		ExposeHeaders:    []string{"X-CSRF-Token"},
		AllowCredentials: true,
	}))
	h.e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper: func(c *echo.Context) bool {
			// skip the csrf for /webhook path
			return c.Path() == "/api/v1/hook/slack"
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

func (h *HttpComponent) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			// 1. Check for session cookie
			cookie, err := ctx.Cookie(utils.SESSION_COOKIE_NAME)
			if err != nil {
				// Return 401 Unauthorized if cookie is missing
				return ctx.JSON(http.StatusUnauthorized, utils.NewAuthenticationError("missing session cookie", err))
			}
			// 2. Validate with Redis
			sessionID := cookie.Value
			sessionData, err := h.redisclient.Get("session:" + sessionID)
			if err != nil {
				// Return 401 Unauthorized if session is invalid or expired in Redis
				return ctx.JSON(http.StatusUnauthorized, utils.NewAuthenticationError("invalid or expired session", err))
			}
			// 3. (Optional) Inject the session data into the context so your GetAuthMe handler can use it
			ctx.Set("user_session", sessionData)
			ctx.Set("session_id", "session:"+sessionID)
			// 4. Proceed to the actual handler (GetAuthMe)
			return next(ctx)
		}
	}
}



