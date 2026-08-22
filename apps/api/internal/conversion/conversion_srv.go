package conversion

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
	"gorm.io/datatypes"
)

type ConverstionSrv struct {
	repo *repositories.ConversationRepo
}

func NewConverstionSrv(dbcon *client.PostgresClient) *ConverstionSrv {
	return &ConverstionSrv{
		repo: repositories.NewConversationRepo(dbcon),
	}
}

func (s *ConverstionSrv) GetConversationBySessionId(sessionId datatypes.UUID) ([]*models.Conversation, error) {
	return s.repo.GetConversationBySessionId(sessionId)
}

func (s *ConverstionSrv) CreateConversation(conversation *models.Conversation) error {
	return s.repo.CreateConversation(conversation)
}

