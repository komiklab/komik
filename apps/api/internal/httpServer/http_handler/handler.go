package httphandler

import (
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/event"
)

type HttpHandler struct {
	cfg      *internal.Config
	dbclient *client.PostgresClient
	publisher *event.Publisher
}

func NewHttpHandler(cfg *internal.Config, dbclient *client.PostgresClient, publisher *event.Publisher) *HttpHandler {
	return &HttpHandler{
		cfg:      cfg,
		dbclient: dbclient,
		publisher: publisher,
	}
}

var _ apidefn.ServerInterface = (*HttpHandler)(nil)
