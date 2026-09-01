package repositories

import (
	"github.com/google/uuid"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
)

type ConversationRepo struct {
	dbcon *client.PostgresClient
}

func (r *ConversationRepo) UpdateConversation(conversation *models.Conversation) error {
	return r.dbcon.GetClient().Save(conversation).Error
}

func NewConversationRepo(dbcon *client.PostgresClient) *ConversationRepo {
	return &ConversationRepo{dbcon: dbcon}
}

func (r *ConversationRepo) GetConversationBySessionId(sessionId uuid.UUID) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	err := r.dbcon.GetClient().Where("session_id = ? ORDER BY created_at DESC", sessionId).Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *ConversationRepo) CreateConversation(conversation *models.Conversation) error {
	return r.dbcon.GetClient().Create(conversation).Error
}
