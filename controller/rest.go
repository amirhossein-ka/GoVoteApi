package controller

type Rest interface {
	Start(addr string) error
	Stop() error
}
