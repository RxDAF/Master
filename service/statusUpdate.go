package service

import (
	"errors"
	"log"

	rmaster "github.com/RxDAF/Master/dto"
)

func (s *RService) StatusUpdate(r rmaster.RMaster_StatusUpdateServer) error {
	// 已经建立了连接 开始绑定连接
	info, err := r.Recv()
	if err != nil {
		return err
	}
	certification := info.GetCertification()
	if certification == nil {
		return errors.New("certification should come first")
	}
	s.serversLock.Lock()
	statusUpdater := newStatusUpdater(r)
	server := s.servers[certification.Address]
	server.StatusUpdater = statusUpdater
	for {
		// 检测是否有消息
		msg, err := r.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		statusUpdate := msg.GetService()
		if statusUpdate == nil {
			log.Println("server:", certification.Address, " statusUpdate error")
			break
		}
		s.servers[certification.Address].UpdateServiceOnlineStatus(statusUpdate.ServiceName, statusUpdate.NewStatus)
	}
	// 开始善后 清除该服务器
	statusUpdater.closed = true
	s.serversLock.Lock()
	defer s.serversLock.Unlock()
	for serviceName := range server.Services {
		// 从该service的server列表里清除掉当前server
		service := s.services[serviceName]
		func() {
			service.Lock.Lock()
			defer service.Lock.Unlock()
			found := -1
			for index, s := range service.MountedServer {
				if s == server {
					found = index
					break
				}
			}
			if found == -1 { //理论上不会有这种情况
				return
			}
			// 开始删除
			service.MountedServer = append(service.MountedServer[:found], service.MountedServer[found+1:]...)
		}()
	}
	return nil
}

type statusUpdater struct {
	conn   rmaster.RMaster_StatusUpdateServer
	closed bool
}

func newStatusUpdater(conn rmaster.RMaster_StatusUpdateServer) *statusUpdater {
	return &statusUpdater{
		conn: conn,
	}
}

// SetServiceStatus 令服务上下线
func (s *statusUpdater) SetServiceStatus(serviceName string, status bool) error {
	if s.closed {
		return errors.New("the connection has been closed")
	}
	return s.conn.Send(&rmaster.StatusUpdateReq{
		StatusUpdate: &rmaster.StatusUpdateReq_Service{
			Service: &rmaster.ServiceStatusChange{
				ServiceName: serviceName,
				NewStatus:   status,
			},
		},
	})
}
