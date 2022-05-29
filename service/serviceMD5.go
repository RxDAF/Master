package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"

	rmaster "github.com/RxDAF/Master/dto"
)

func (s *RService) loadServiceFile(serviceName string) (io.ReadCloser, error) {
	path := s.cfg.ServicesPath + string(os.PathSeparator) + serviceName + ".tar"
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (s *RService) ServiceFileMD5(ctx context.Context, r *rmaster.ServiceFileMD5Request) (*rmaster.ServiceFileMD5Result, error) {
	var err error
	md5str, ok := func() (string, bool) {
		s.serviceMD5RecorderLock.RLock()
		defer s.serviceMD5RecorderLock.RUnlock()
		md5, ok := s.serviceMD5Recorder[r.ServiceName]
		return md5, ok
	}()
	if !ok {
		// 那就要读取
		s.serviceMD5RecorderLock.Lock()
		defer s.serviceMD5RecorderLock.Unlock()
		var res io.ReadCloser
		res, err = s.loadServiceFile(r.ServiceName)
		defer res.Close()
		md5Handle := md5.New()
		_, err = io.Copy(md5Handle, res)
		if nil != err {
			return nil, err
		}
		md := md5Handle.Sum(nil)
		md5str = fmt.Sprintf("%x", md)
	}
	return &rmaster.ServiceFileMD5Result{Md5: md5str}, err
}
