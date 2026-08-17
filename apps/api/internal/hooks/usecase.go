package hooks

import (
	"context"

	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/entity"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
)

const (
	SLACK_HOOK_NAME string = "slack"
)

type HookService struct {
	hookRepo  *repositories.HooksRepo
	entitySrv *entity.EntityService
}

func NewHookService(dbclient *client.PostgresClient) *HookService {
	return &HookService{
		hookRepo:  repositories.NewHooksRepo(dbclient),
		entitySrv: entity.NewEntityService(dbclient),
	}
}

func (h *HookService) FetchHooks() ([]models.Hooks, error) {
	return h.hookRepo.FetchHooks()
}

func (h *HookService) CreateHook(hook *models.Hooks) error {
	if hook.Name == "" {
		return nil
	}
	return h.hookRepo.CreateHook(hook)
}

func (h *HookService) SendMessage(ctx context.Context, hookinput *models.HookInput) error {
	return h.entitySrv.InitiateEntity(ctx, hookinput.Name, hookinput.Reference, hookinput.Initiator, hookinput.Message)
}

func (h *HookService) SupportedHookName(ctx context.Context) []string {
	custom_hooks, err := h.hookRepo.FetchHooks()
	if err != nil {
		return nil
	}
	var hookNames []string
	hookNames = append(hookNames, SLACK_HOOK_NAME)
	for _, hook := range custom_hooks {
		hookNames = append(hookNames, hook.Name)
	}
	return hookNames
}
