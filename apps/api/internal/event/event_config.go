package event

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/nats-io/nats.go/jetstream"
	js "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/jetstream"
)

const (
	StreamName = "komik-stream"
	AuditLogSubject = "komik.>"
	AuditLogHandlerName = "komik-audit-handler"
)


func consumerForSubject(subject string) js.ResourceInitializer {
	return func(ctx context.Context, jsCtx jetstream.JetStream, topic string) (jetstream.Consumer, func(context.Context, watermill.LoggerAdapter), error) {
		stream, err := jsCtx.Stream(ctx, StreamName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get stream for topic %s: %w", topic, err)
		}
		consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
			FilterSubject: subject,
			AckPolicy:    jetstream.AckExplicitPolicy,
			DeliverPolicy: jetstream.DeliverAllPolicy,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create or update consumer for topic %s: %w", topic, err)
		}
		return consumer, nil, nil
	}
}