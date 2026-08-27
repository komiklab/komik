package repositories

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
	"gorm.io/datatypes"
)

type AgentRepo struct {
	dbcon *client.PostgresClient
}

func (a *AgentRepo) ListAgentsBasedOnCustomHook(entity models.Entity) ([]models.Agent, error) {
	gormdb := a.dbcon.GetClient()
	hook := entity.SourceType
	var agents []models.Agent
	result := gormdb.Where(datatypes.JSONQuery("triggered_by").HasKey(hook)).Find(&agents)
	if result.Error != nil {
		return nil, utils.NewGeneralError(result.Error)
	}
	return agents, nil

}

func NewAgentRepo(dbcon *client.PostgresClient) *AgentRepo {
	return &AgentRepo{
		dbcon: dbcon,
	}
}

func (a *AgentRepo) ListAgents() ([]models.Agent, error) {
	gormdb := a.dbcon.GetClient()
	var agents []models.Agent
	err := gormdb.Find(&agents).Error
	if err != nil {
		return nil, utils.NewGeneralError(err)
	}
	return agents, nil
}

func (a *AgentRepo) CreateAgent(agent *models.Agent) error {
	gormdb := a.dbcon.GetClient()
	err := gormdb.Create(agent).Error
	if err != nil {
		return utils.NewGeneralError(err)
	}
	return nil
}
