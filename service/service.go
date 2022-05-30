package service

import (
	"sync"

	"github.com/RxDAF/Master/cfg"
)

type RService struct {
	cfg                    *cfg.Configure
	services               map[string]*service
	servicesLock           sync.RWMutex
	servers                map[string]*server
	serversLock            sync.RWMutex
	serviceMD5Recorder     map[string]string
	serviceMD5RecorderLock sync.RWMutex
}
type service struct {
	MountedServer []*server
	Lock          sync.RWMutex
}
type server struct {
	Address       string
	Services      map[string]*serviceStatus
	servicesLock  sync.RWMutex
	StatusUpdater *statusUpdater
}
type serviceStatus struct {
	Online bool //是否在线
}

func (s *server) UpdateServiceOnlineStatus(serviceName string, Online bool) {
	s.servicesLock.RLock()
	defer s.servicesLock.RUnlock()
	s.Services[serviceName].Online = Online
}
func NewService(cfg *cfg.Configure) *RService {
	return &RService{
		cfg:                cfg,
		services:           make(map[string]*service),
		serviceMD5Recorder: make(map[string]string),
		servers:            make(map[string]*server),
	}
}
