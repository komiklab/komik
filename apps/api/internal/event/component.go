package event

import (
	"context"

	"github.com/komiklab/komik/internal/component"
)

type WatermillComponent struct {
}

// GetName implements [Component].
func (w *WatermillComponent) GetName() string {
	return "WatermillComponent"
}

// Init implements [Component].
func (w *WatermillComponent) Init() {
	panic("unimplemented")
}

// Start implements [Component].
func (w *WatermillComponent) Start() {
	panic("unimplemented")
}

// Stop implements [Component].
func (w *WatermillComponent) Stop(ctx context.Context) {
	panic("unimplemented")
}

var _ component.Component = (*WatermillComponent)(nil)
