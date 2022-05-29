package service

import (
	"context"

	rmaster "github.com/RxDAF/Master/dto"
)

func (s *RService) MountServer(serviceName string, serverInformation *server) {
	rservice := func() *service {
		s.servicesLock.Lock()
		defer s.servicesLock.Unlock()
		ret, ok := s.services[serviceName]
		if !ok {
			ret = &service{
				MountedServer: make([]*server, 0),
			}
			s.services[serviceName] = ret
		}
		return ret
	}()
	rservice.Lock.Lock()
	defer rservice.Lock.Unlock()
	rservice.MountedServer = append(rservice.MountedServer, serverInformation)
}
func (s *RService) RegisterService(ctx context.Context, r *rmaster.RegisterServiceRequest) (*rmaster.RegisterServiceResult, error) {
	serverInformation := &server{
		Address: r.Address,
	}
	for _, serviceName := range r.Roles {
		s.MountServer(serviceName, serverInformation)
	}
	return &rmaster.RegisterServiceResult{}, nil
}
