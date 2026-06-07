package component

import (
	"context"
	"net/http"
	"time"

	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal"
	httphandler "github.com/komiklab/komik/internal/handler/http_handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog/log"
)

type HttpComponent struct {
	e *echo.Echo
	s *http.Server
	cfg *internal.Config
}

var _ Component = (*HttpComponent)(nil)

func (h *HttpComponent) GetName() string {
	return "HttpComponent"
}

func (h *HttpComponent) Init() {
	h.e = echo.New()
	h.addMiddleware()
	handler := httphandler.NewHttpHandler()
	apidefn.RegisterHandlersWithOptions(h.e, handler, apidefn.RegisterHandlersOptions{
		BaseURL:              "/api/v1",
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{},
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
	if h.s != nil {
		if err := h.s.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("HTTP server failed to shutdown")
		}
	}
}

func NewHttpComponent(cfg *internal.Config) *HttpComponent {
	return &HttpComponent{
		cfg: cfg,
	}
}

func (h *HttpComponent) addMiddleware() {
	h.e.Use(middleware.Recover())
	h.e.Use(middleware.RequestID())
	h.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://example.com"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	h.e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieSecure:   true,
		CookieHTTPOnly: false,
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
