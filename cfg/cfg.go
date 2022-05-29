package cfg

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Configure struct {
	ServicesPath string `yaml:"servicesPath"` // 各个微服务的压缩包文件目录
}

func NewCFG(file string) (*Configure, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	cfg := new(Configure)
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
