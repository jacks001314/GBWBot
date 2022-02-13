package client

type Client interface {

	Start() error

	Stop()
}
