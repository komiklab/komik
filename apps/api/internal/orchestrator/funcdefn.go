package orchestrator

import (

	"github.com/rs/zerolog/log"
)

func (o *Orchestrator) RegisterFuncs() {
	log.Debug().Msg("registering functions")
	// register all the functions
	o.registerEntityTransitionFn()

}
