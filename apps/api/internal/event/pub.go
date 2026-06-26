package event

import (
	js "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/jetstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/komiklab/komik/internal"
	"github.com/rs/zerolog/log"
	"github.com/alexdrl/zerowater"
)

type Publisher struct {
	publisher *js.Publisher
}

func NewPublisher(cfg *internal.Config) *Publisher {
	logger := zerowater.NewZerologLoggerAdapter(log.Logger)
	publisher, err := js.NewPublisher(js.PublisherConfig{
		Logger: logger,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create publisher")
	}
	return &Publisher{
		publisher: publisher,
	}
}

func (p *Publisher) GetPublisher() *js.Publisher {
	return p.publisher
}

func (p *Publisher) Publish(subject string, message *message.Message) error {
	log.Debug().Str("subject", subject).Msg("Publishing event")
	return p.publisher.Publish(subject, message)
}
