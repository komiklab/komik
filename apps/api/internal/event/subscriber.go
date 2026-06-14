package event

import (
	"fmt"
	"time"

	js "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/jetstream"
	"github.com/komiklab/komik/internal"
	"github.com/rs/zerolog/log"
)

type NatsSubscriber struct {
	inner *js.Subscriber
}

func NewNatsSubscriber(cfg *internal.Config, subject string) (*NatsSubscriber, error) {
	logger := NewZerologLoggerAdapter(log.Logger)
	subscriber, err := js.NewSubscriber(js.SubscriberConfig{
		URL:                 cfg.NatsURL,
		Logger:              logger,
		AckWaitTimeout:      30 * time.Second, // 30 seconds
		ResourceInitializer: consumerForSubject(subject),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create nats subscriber: %w", err)
	}
	return &NatsSubscriber{
		inner: subscriber,
	}, nil
}

func (s *NatsSubscriber) GetSubscriber() *js.Subscriber {
	return s.inner
}
