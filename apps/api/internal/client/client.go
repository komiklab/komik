package client

type Client interface {
	GetClient() any
	Ping() error
	Close()
}
