package rmaster

import (
	"bytes"
	"context"

	grpc "google.golang.org/grpc"
)

type RMaster struct {
	client RMasterClient
}

func NewRMaster(serviceAddress string) (*RMaster, error) {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure()) // 建立链接
	if err != nil {
		return nil, err
	}
	client := NewRMasterClient(conn) // 初始化客户端
	return &RMaster{
		client: client,
	}, nil
}
func (r *RMaster) RegisterService(address string, roles []string) error {
	_, err := r.client.RegisterService(context.Background(), &RegisterServiceRequest{
		Address: address,
		Roles:   roles,
	})
	return err
}
func (r *RMaster) ServiceFileMD5(serviceName string) (md5 string, err error) {
	var res *ServiceFileMD5Result
	res, err = r.client.ServiceFileMD5(context.Background(), &ServiceFileMD5Request{
		ServiceName: serviceName,
	})
	if err != nil {
		return
	}
	md5 = res.Md5
	return
}
func (r *RMaster) DownloadService(serviceName string) ([]byte, error) {
	res, err := r.client.DownloadService(context.Background(), &DownloadServiceRequest{
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	for {
		resp, err := res.Recv()
		if err != nil {
			return nil, err
		}
		if len(resp.Data) == 0 { // 下载完成
			break
		}
		_, err = buf.Write(resp.Data)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func (r *RMaster) SetupConnect() (RMaster_StatusUpdateClient, error) {
	conn, err := r.client.StatusUpdate(context.Background())
	return conn, err
}
