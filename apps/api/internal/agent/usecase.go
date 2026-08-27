package agent

import (
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
)

type AgentService struct {
	agentRepo *repositories.AgentRepo
}

func (a *AgentService) CreateAgent(agentModel *models.Agent) (*apidefn.AgentCreateResponse, error) {
	err := a.agentRepo.CreateAgent(agentModel)
	if err != nil {
		return nil, err
	}
	return &apidefn.AgentCreateResponse{
		AgentID: agentModel.Id,
		Name:    agentModel.Name,
	}, nil
}

func NewAgentService(dbcon *client.PostgresClient) *AgentService {
	agentRepo := repositories.NewAgentRepo(dbcon)
	return &AgentService{
		agentRepo: agentRepo,
	}
}

func (a *AgentService) ListAgent() (*models.ListAgents, error) {
	agents, err := a.agentRepo.ListAgents()
	if err != nil {
		return nil, err
	}
	return &models.ListAgents{
		Agents: agents,
	}, nil
}

func (a *AgentService) ListAgentBasedOnEntity(entity *models.Entity) (*models.ListAgents, error) {
	agents, err := a.agentRepo.ListAgentsBasedOnEntity(entity)
	if err != nil {
		return nil, err
	}
	return &models.ListAgents{
		Agents: agents,
	}, nil
}
