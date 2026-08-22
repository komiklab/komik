package conversion

import (
	"context"

	"github.com/komiklab/komik/internal/models"
	"gorm.io/datatypes"
)

func FetchConversationBySessionId(ctx context.Context, srv *ConverstionSrv, sessionId datatypes.UUID) ([]*models.Conversation, error) {
	return srv.repo.GetConversationBySessionId(sessionId)
}

func CreateConversation(ctx context.Context, srv *ConverstionSrv, conversation *models.Conversation) error {
	return srv.repo.CreateConversation(conversation)
}