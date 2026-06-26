package repositories

import (
	"database/sql"
	"errors"

	sqlpubsub "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alexdrl/zerowater"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ForwarderPublisher struct {
	logger *zerowater.ZerologLoggerAdapter
}

func NewForwarderPublisher() *ForwarderPublisher {
	logger := zerowater.NewZerologLoggerAdapter(log.Logger)
	return &ForwarderPublisher{
		logger: logger,
	}
}

func (f *ForwarderPublisher) Publish(tx *gorm.DB, topic string, msg *message.Message) error {
	sqlTx, ok := tx.Statement.ConnPool.(*sql.Tx)
	if !ok {
		return errors.New("gorm transaction is not a sql transaction")
	}
	sqlpub, err := sqlpubsub.NewPublisher(
		sqlpubsub.TxFromStdSQL(sqlTx),
		sqlpubsub.PublisherConfig{
			SchemaAdapter: sqlpubsub.DefaultPostgreSQLSchema{},
		},
		f.logger,
	)
	if err != nil {
		return err
	}
	fwdpub := forwarder.NewPublisher(sqlpub, forwarder.PublisherConfig{})
	err = fwdpub.Publish(topic, msg)
	if err != nil {
		return err
	}
	return nil
}
