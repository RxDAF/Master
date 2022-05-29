package service

import (
	"sync"

	"github.com/RxDAF/Master/cfg"
)

type RService struct {
	cfg                    *cfg.Configure
	services               map[string]*service
	servicesLock           sync.RWMutex
	serviceMD5Recorder     map[string]string
	serviceMD5RecorderLock sync.RWMutex
}
type service struct {
	MountedServer []*server
	Lock          sync.RWMutex
}
type server struct {
	Address string
}

func NewService(cfg *cfg.Configure) *RService {
	return &RService{
		cfg:                cfg,
		services:           make(map[string]*service),
		serviceMD5Recorder: make(map[string]string),
	}
}
