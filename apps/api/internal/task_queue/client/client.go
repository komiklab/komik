package client

import (
	"github.com/hibiken/asynq"
	"github.com/komiklab/komik/internal/client"
	taskqueue "github.com/komiklab/komik/internal/task_queue"
)

type AsynqClient struct {
	client *asynq.Client
}

func NewAsyncClient(redisClient *client.RedisClient) *AsynqClient {
	return &AsynqClient{
		client: asynq.NewClientFromRedisClient(redisClient.GetClient()),
	}
}

func (a *AsynqClient) EnqueueTask(taskType taskqueue.TaskType, payload []byte, headers map[string]string) error {
	// ensure taskType is valid
	_, err := taskqueue.ParseTaskType(taskType)
	if err != nil {
		return err
	}
	_, err = a.client.Enqueue(asynq.NewTaskWithHeaders(string(taskType), payload, headers))
	return err
}

func (a *AsynqClient) GetClient() *asynq.Client {
	return a.client
}
