package agent

import (
	"bufio"
	"encoding/json"
	"strings"

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
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "text/event-stream, application/json").
		SetBody(body).
		SetResponseDoNotParse(true).
		Post(url)
	if err != nil {
		return nil, err
	}
	defer res.RawResponse.Body.Close()
	contentType := res.Header().Get("Content-Type")
	var events []json.RawMessage
	if strings.Contains(contentType, "text/event-stream"){
		scanner := bufio.NewScanner(res.RawResponse.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(line[5:])
				if data == "" || data == "[DONE]" {
					continue
				}
				events = append(events, json.RawMessage(data))
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return json.Marshal(events)
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
