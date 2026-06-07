package component

import (
	"context"
)

type Component interface {
	GetName() string
	Init()
	Start()
	Stop(ctx context.Context)
}
