package httphandler

import (
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
)

type HttpHandler struct {
	cfg      *internal.Config
	dbclient client.Client
}

func NewHttpHandler(cfg *internal.Config, dbclient client.Client) *HttpHandler {
	return &HttpHandler{
		cfg:      cfg,
		dbclient: dbclient,
	}
}

var _ apidefn.ServerInterface = (*HttpHandler)(nil)
