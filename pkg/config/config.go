package config

import (
	"encoding/json"
	"github.com/containernetworking/cni/pkg/types"
)

// Config CNI conf配置文件
type Config struct {
	//types.NetConf // 内置一些字段，直接嵌套使用，弃用，缺少一些ipam需要的字段，自己写
	CNIVersion string `json:"cniVersion"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Bridge     string `json:"bridge"`
	IPAM       IPAM   `json:"ipam"`
}

type IPAM struct {
	Type       string         `json:"type"`
	Subnet     string         `json:"subnet"`
	RangeStart string         `json:"rangeStart"`
	RangeEnd   string         `json:"rangeEnd"`
	DataDir    string         `json:"dataDir"`
	Routes     []*types.Route `json:"routes"`
}

func LoadCNIConfig(data []byte) (*Config, error) {

	cfg := &Config{}
	err := json.Unmarshal(data, cfg)
	if err != nil {

		return nil, err
	}
	return cfg, nil
}
