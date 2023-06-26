package main

import (
	"fmt"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	current "github.com/containernetworking/cni/pkg/types/040"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ipam"
	"github.com/vishvananda/netlink"
	"github/mycni/cni_practice/pkg/bridge"
	"github/mycni/cni_practice/pkg/config"
	"github/mycni/cni_practice/pkg/veth"
)

func log(str string) {
	fmt.Fprintf(os.Stderr, str)
}

func cmdAdd(args *skel.CmdArgs) error {
	cfg, err := config.ConfigFromStdin(args.StdinData)
	if err != nil {
		return err
	}

	ret := &current.Result{CNIVersion: cfg.CNIVersion}

	if cfg.IPAM.Type != "" {
		r, err := ipam.ExecAdd(cfg.IPAM.Type, args.StdinData)
		if err != nil {
			return err
		}
		ipamRet, err := current.NewResultFromResult(r)
		if err != nil {
			return err
		}
		//到这一步获取到 分配的IP
		ret.IPs = ipamRet.IPs
		ret.DNS = ipamRet.DNS
		ret.Routes = ipamRet.Routes
		//ipamRet.PrintTo(os.Stderr)

	}

	// 创建或更新网桥
	var br *netlink.Bridge
	if br, err = bridge.CreateOrUpdateBridge("jtthink0"); err != nil {
		return err
	}

	// 创建veth 设备
	err = veth.CreateVeth(args.Netns, ret.IPs[0].Address.String(), br)
	if err != nil {
		return err
	}

	return ret.Print()

}

func main() {
	skel.PluginMain(cmdAdd, nil, nil, version.All, "jtthink")
}
