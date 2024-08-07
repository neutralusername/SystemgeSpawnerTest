package appWebsocketHTTP

import "sync/atomic"

type AppWebsocketHTTP struct {
	ports *atomic.Uint32
}

func New() *AppWebsocketHTTP {
	ports := &atomic.Uint32{}
	ports.Store(32768)
	return &AppWebsocketHTTP{
		ports: ports,
	}
}
