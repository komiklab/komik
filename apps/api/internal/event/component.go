package event

import (
	"context"
	"database/sql"
	"time"

	sqlpubsub "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/alexdrl/zerowater"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	"github.com/rs/zerolog/log"
)

type EventLoopComponent struct {
	router         *message.Router
	cfg            *internal.Config
	dbcon          *client.PostgresClient
	temporalClient *client.TemporalClient
}

func NewEventLoopComponent(cfg *internal.Config, dbcon *client.PostgresClient, temporalClient *client.TemporalClient) *EventLoopComponent {
	return &EventLoopComponent{
		cfg:            cfg,
		dbcon:          dbcon,
		temporalClient: temporalClient,
	}
}

// GetName implements [Component].
func (w *EventLoopComponent) GetName() string {
	return "EventLoopComponent"
}

// Init implements [Component].
func (w *EventLoopComponent) Init() {
	logger := zerowater.NewZerologLoggerAdapter(log.Logger)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create router")
	}
	w.router = router
	w.router.AddPlugin(plugin.SignalsHandler)
	w.router.AddMiddleware(middleware.CorrelationID, middleware.Recoverer, middleware.Retry{
		MaxRetries:      3,
		InitialInterval: 100 * time.Millisecond,
		Logger:          logger,
	}.Middleware)
	auditSub, err := NewNatsSubscriber(w.cfg, AuditLogSubject)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create audit subscriber")
	}
	entityInitiatedSub, err := NewNatsSubscriber(w.cfg, EntitySubjectInitiated)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create entity subscriber")
	}
	handler := NewEventHandler(w.dbcon, w.temporalClient)
	w.router.AddConsumerHandler(AuditLogHandlerName, StreamName, auditSub.GetSubscriber(), handler.HandleAuditLog)
	w.router.AddConsumerHandler(EntityHandlerName, StreamName, entityInitiatedSub.GetSubscriber(), handler.HandleEntityInitiated)
	w.forwarder(logger)
}

// Start implements [Component].
func (w *EventLoopComponent) Start() {
	err := w.router.Run(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start router")
	}
	log.Debug().Msg("Router started")
}

// Stop implements [Component].
func (w *EventLoopComponent) Stop(ctx context.Context) {
	err := w.router.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to stop router")
	}
	log.Debug().Msg("Router stopped")
}

func (w *EventLoopComponent) forwarder(logger *zerowater.ZerologLoggerAdapter) {
	dsn := w.cfg.PostgresDSN
	stdDb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	sqlSubscriber, err := sqlpubsub.NewSubscriber(
		sqlpubsub.BeginnerFromStdSQL(stdDb), sqlpubsub.SubscriberConfig{
			SchemaAdapter:    sqlpubsub.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   sqlpubsub.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: true,
		}, logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create SQL subscriber")
	}
	pub := NewPublisher(w.cfg)
	_, err = forwarder.NewForwarder(sqlSubscriber, pub.GetPublisher(), logger, forwarder.Config{
		Router: w.router,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create forwarder")
	}
}

var _ internal.Component = (*EventLoopComponent)(nil)
