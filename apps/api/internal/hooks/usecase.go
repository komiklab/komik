package hooks

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
)

type HookService struct {
	hookRepo *repositories.HooksRepo
}

func NewHookService(dbclient *client.PostgresClient) *HookService {
	return &HookService{
		hookRepo: repositories.NewHooksRepo(dbclient),
	}
}

func (h *HookService) FetchHooks() ([]models.Hooks, error) {
	return h.hookRepo.FetchHooks()
}

func (h *HookService) CreateHook(hook *models.Hooks) error {
	if hook.Name == ""{
		return nil
	}
	return h.hookRepo.CreateHook(hook)
}