package httphandler

import (
	"net/http"

	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/agent"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
	"github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"
)

// GetAgent implements [apidefn.ServerInterface].
func (h *HttpHandler) GetAgent(ctx *echo.Context) error {
	svc := agent.NewAgentService(h.dbclient)
	agents, err := svc.ListAgent()
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}
	return ctx.JSON(http.StatusOK, agents)
}

func (h *HttpHandler) InternalGetAgents(ctx *echo.Context) error {
	return h.GetAgent(ctx)
}

// GetAgentId implements [apidefn.ServerInterface].
func (h *HttpHandler) GetAgentId(ctx *echo.Context, id types.UUID) error {
	panic("unimplemented")
}

// PostAgent implements [apidefn.ServerInterface].
func (h *HttpHandler) PostAgent(ctx *echo.Context) error {
	log.Info().Msg("called PostAgent")
	var req apidefn.AgentCreateRequest
	if err := ctx.Bind(&req); utils.IsErrNotNil(err) {
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to bind request", err))
	}
	agentModel,err:=models.NewAgentFromApiDefn(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create agent model")
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to create agent model", err))
	}
	svc := agent.NewAgentService(h.dbclient)
	resp, err := svc.CreateAgent(agentModel)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.NewBindError("failed to create agent model", err))
	}
	//TODO: Fire an event[KOM-21]
	return ctx.JSON(http.StatusOK, resp)
}
