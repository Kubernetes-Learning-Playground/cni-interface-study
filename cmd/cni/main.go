package main

import (
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	current "github.com/containernetworking/cni/pkg/types/040"
	"github.com/containernetworking/cni/pkg/version"
	"golanglearning/new_project/cni_practice/pkg/bridge"
	"golanglearning/new_project/cni_practice/pkg/config"
	"github.com/containernetworking/plugins/pkg/ipam"

	"os"
)

func log(str string) {
	fmt.Fprintf(os.Stderr, str)
}

func cmdAdd(args *skel.CmdArgs) error {
	cfg, err := config.ConfigFromStdin(args.StdinData)
	if err != nil {
		return err
	}

	if cfg.IPAM.Type != "" {
		r, err := ipam.ExecAdd(cfg.IPAM.Type, args.StdinData)
		if err != nil {
			return err
		}
		ipamRet, err := current.NewResultFromResult(r)
		if err != nil {
			return err
		}
		ipamRet.PrintTo(os.Stderr)

	}

	// 创建或更新网桥
	if _, err = bridge.CreateOrUpdateBridge("jtthink0"); err != nil {
		return err
	}

	ret := &current.Result{CNIVersion: cfg.CNIVersion}




	return ret.Print()

}

func main() {
	skel.PluginMain(cmdAdd, nil, nil, version.All, "jtthink")
}
