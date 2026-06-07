package client

import (
	"github.com/komiklab/komik/internal"
)

var _ Client = (*PostgresClient)(nil)

type PostgresClient struct {
}

func (p *PostgresClient) GetClient() any {
	// TODO implement me
	panic("implement me")
}

func (p *PostgresClient) Ping() error {
	// TODO implement me
	return nil
}

func (p *PostgresClient) Close() {
	// TODO implement me
	panic("implement me")
}

func NewPostgresClient(cfg *internal.Config) *PostgresClient {
	return &PostgresClient{}
}
