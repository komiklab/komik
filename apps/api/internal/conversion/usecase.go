package conversion

import (
	"context"

	"github.com/komiklab/komik/internal/models"
	"github.com/google/uuid"
)

func FetchConversationBySessionId(ctx context.Context, srv *ConverstionSrv, sessionId uuid.UUID) ([]*models.Conversation, error) {
	return srv.repo.GetConversationBySessionId(sessionId)
}

func CreateConversation(ctx context.Context, srv *ConverstionSrv, conversation *models.Conversation) error {
	return srv.repo.CreateConversation(conversation)
}