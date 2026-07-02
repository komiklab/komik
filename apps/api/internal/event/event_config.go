package event

import (
	"context"
	"fmt"
	"strings"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/nats-io/nats.go/jetstream"
	js "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/jetstream"
)

const (
	StreamName = "komik-stream"
	AuditLogSubject = "komik.>"
	AuditLogHandlerName = "komik-audit-handler"
	EntitySubject = "komik.entity.>"
	EntitySubjectDispatched = "komik.entity.dispatched"
	EntitySubjectInitiated = "komik.entity.initiated"
	EntityHandlerName = "komik-entity-handler"
)


func consumerForSubject(subject string) js.ResourceInitializer {
	return func(ctx context.Context, jsCtx jetstream.JetStream, topic string) (jetstream.Consumer, func(context.Context, watermill.LoggerAdapter), error) {
		stream, err := jsCtx.Stream(ctx, StreamName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get stream for topic %s: %w", topic, err)
		}
		// sanitize: replace dots and strip wildcards
        durableName := strings.NewReplacer(
            ".", "-",
            ">", "",
            "*", "",
        ).Replace(subject)
		durableName = strings.Trim(durableName, "-")
		consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:         durableName,
			FilterSubject: subject,
			AckPolicy:    jetstream.AckExplicitPolicy,
			DeliverPolicy: jetstream.DeliverNewPolicy,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create or update consumer for topic %s: %w", topic, err)
		}
		return consumer, nil, nil
	}
}