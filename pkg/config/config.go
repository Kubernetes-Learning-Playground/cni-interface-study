package config

import (
	"encoding/json"

	"github.com/containernetworking/cni/pkg/types"
)

type Config struct {
	// 内置一些字段，直接嵌套使用
	types.NetConf
	Bridge string `json:"bridge"` // 完全名称。 默认是 jtthink0
}

func ConfigFromStdin(data []byte) (*Config, error) {

	cfg := &Config{}
	err := json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
