package repositories

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"gorm.io/datatypes"
)

type ConversationRepo struct {
	dbcon *client.PostgresClient
}



func NewConversationRepo(dbcon *client.PostgresClient) *ConversationRepo {
	return &ConversationRepo{dbcon: dbcon}
}

func (r *ConversationRepo) GetConversationBySessionId(sessionId datatypes.UUID) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	err := r.dbcon.GetClient().Where("session_id = ?", sessionId).Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *ConversationRepo) CreateConversation(conversation *models.Conversation) error {
	return r.dbcon.GetClient().Create(conversation).Error
}
