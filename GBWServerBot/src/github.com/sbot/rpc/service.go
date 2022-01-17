package rpc

type Service interface {

	Start()
	Stop()

	Name() string

}
