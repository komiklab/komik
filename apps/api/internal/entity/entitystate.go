package entity

import (
	"context"

	"github.com/komiklab/komik/internal/models"
	"github.com/looplab/fsm"
)

const (
	EntityStateInitiated  string = "initiated"
	EntityStateDispatched string = "dispatched"
)

const (
	EntityEventDispatch string = "dispatch"
)

type EntityTransitioner struct {
	fsm    *fsm.FSM
	entity *models.Entity
}

func (e *EntityTransitioner) Dispatch(ctx context.Context) error {
	return e.fsm.Event(ctx, EntityEventDispatch)
}

func NewEntityTransitioner(entity *models.Entity) *EntityTransitioner {
	et := EntityTransitioner{
		entity: entity,
	}
	et.fsm = fsm.NewFSM(
		entity.Status,
		fsm.Events{
			{Name: EntityEventDispatch, Src: []string{EntityStateInitiated}, Dst: EntityStateDispatched},
		},
		fsm.Callbacks{},
	)
	return &et
}
