package httphandler

import (
	"github.com/komiklab/komik/apidefn"
)

type HttpHandler struct {
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

var _ apidefn.ServerInterface = (*HttpHandler)(nil)
