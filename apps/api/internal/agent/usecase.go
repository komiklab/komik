package agent

import (
	"encoding/json"

	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
	"resty.dev/v3"
)

type AgentService struct {
	agentRepo *repositories.AgentRepo
}

func (a *AgentService) CallAgent(agent models.Agent, entity models.Entity) ([]byte, error){
	url := agent.Endpoint
	payloads_required := agent.Parameter
	entitySourcePayload := make(map[string]interface{})
	if err := json.Unmarshal(entity.SourcePayload, &entitySourcePayload); err != nil {
		return  nil, err
	}
	body := map[string]interface{}{}
	for _, param := range payloads_required {
		if val, exists := entitySourcePayload[param.Name]; exists {
			body[param.Name] = val
		}
	}
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	return  nil, err
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

func (a *AgentService) ListAgentBasedOnEntity(entity models.Entity) (*models.ListAgents, error) {
	agents, err := a.agentRepo.ListAgentsBasedOnCustomHook(entity)
	if err != nil {
		return nil, err
	}
	return &models.ListAgents{
		Agents: agents,
	}, nil
}
