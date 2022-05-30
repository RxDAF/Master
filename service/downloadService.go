package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	rmaster "github.com/RxDAF/Master/dto"
)

func (s *RService) DownloadService(req *rmaster.DownloadServiceRequest, res rmaster.RMaster_DownloadServiceServer) error {
	// 判断对应服务是否存在
	_, ok := s.cfg.Services[req.ServiceName]
	if !ok {
		return errors.New("no such service")
	}
	// 开始下载
	filePath := filepath.Join(s.cfg.ServicesPath, req.ServiceName)
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF { // 已经读完了
				return res.Send(&rmaster.DownloadServiceResult{
					Data: buf[0:0],
				})
			}
			return err
		}
		if err := res.Send(&rmaster.DownloadServiceResult{
			Data: buf[0:n],
		}); err != nil {
			return err
		}
	}
}
