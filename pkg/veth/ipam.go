package veth

import (
	"encoding/json"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/plugins/pkg/ipam"
	"github.com/containernetworking/plugins/plugins/ipam/host-local/backend/allocator"
	"github.com/pkg/errors"
	"github/mycni/cni_practice/pkg/config"
	"net"
)

// Ipam 使用ipam插件分配网络地址
func Ipam(conf *config.Config) (types.Result, error) {
	ipNet, err := types.ParseCIDR(conf.IPAM.Subnet)
	if err != nil {
		return nil, errors.Wrapf(err, "ParseCIDR error")
	}

	var startIP, endIP net.IP
	if conf.IPAM.RangeStart != "" {
		startIP = net.ParseIP(conf.IPAM.RangeStart)
		if startIP == nil {
			return nil, errors.Errorf("range start %s error", conf.IPAM.RangeStart)
		}
	}
	if conf.IPAM.RangeEnd != "" {
		endIP = net.ParseIP(conf.IPAM.RangeEnd)
		if endIP == nil {
			return nil, errors.Errorf("range end %s error", conf.IPAM.RangeEnd)
		}
	}

	ipamConf := allocator.Net{
		Name:       conf.Name,
		CNIVersion: conf.CNIVersion,
		IPAM: &allocator.IPAMConfig{
			Type: conf.IPAM.Type,
			Ranges: []allocator.RangeSet{
				{
					{
						Subnet:     types.IPNet(*ipNet),
						RangeStart: startIP,
						RangeEnd:   endIP,
					},
				},
			},
			DataDir: conf.IPAM.DataDir, // 加入才会生效
			Routes:  conf.IPAM.Routes,
		},
	}
	ipamConfBytes, err := json.Marshal(ipamConf)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal ipam conf error")
	}

	return ipam.ExecAdd(conf.IPAM.Type, ipamConfBytes)
}

// ReleaseIP 释放ipam分配的地址
func ReleaseIP(conf *config.Config) error {
	ipNet, err := types.ParseCIDR(conf.IPAM.Subnet)
	if err != nil {
		return errors.Wrapf(err, "ParseCIDR error")
	}

	var startIP, endIP net.IP
	if conf.IPAM.RangeStart != "" {
		startIP = net.ParseIP(conf.IPAM.RangeStart)
		if startIP == nil {
			return errors.Errorf("range start %s error", conf.IPAM.RangeStart)
		}
	}
	if conf.IPAM.RangeEnd != "" {
		endIP = net.ParseIP(conf.IPAM.RangeEnd)
		if endIP == nil {
			return errors.Errorf("range end %s error", conf.IPAM.RangeEnd)
		}
	}

	ipamConf := allocator.Net{
		Name:       conf.Name,
		CNIVersion: conf.CNIVersion,
		IPAM: &allocator.IPAMConfig{
			Type: conf.IPAM.Type,
			Ranges: []allocator.RangeSet{
				{
					{
						Subnet:     types.IPNet(*ipNet),
						RangeStart: startIP,
						RangeEnd:   endIP,
					},
				},
			},
			DataDir: conf.IPAM.DataDir, // 加入才会生效
			Routes:  conf.IPAM.Routes,
		},
	}
	ipamConfBytes, err := json.Marshal(ipamConf)
	if err != nil {
		return errors.Wrapf(err, "marshal ipam conf error")
	}

	return ipam.ExecDel(conf.IPAM.Type, ipamConfBytes)
}
