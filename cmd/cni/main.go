package main

import (
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/040"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
	"github.com/vishvananda/netlink"
	"github/mycni/cni_practice/pkg/bridge"
	"github/mycni/cni_practice/pkg/config"
	"github/mycni/cni_practice/pkg/veth"
)

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("mycniplugin"))
}

// cmdAdd CNI add方法
func cmdAdd(args *skel.CmdArgs) error {
	// 使用cnitool 打印log会报错。
	//fmt.Printf("cmdAdd containerID: %s \n", args.ContainerID)
	//fmt.Printf("cmdAdd netNS: %s \n", args.Netns)
	//fmt.Printf("cmdAdd ifName: %s \n", args.IfName)
	//fmt.Printf("cmdAdd args: %s \n", args.Args)
	//fmt.Printf("cmdAdd path: %s \n", args.Path)
	//fmt.Printf("cmdAdd stdin: %s \n", string(args.StdinData))
	/*
			cmdAdd containerID: cnitool-77383ca0a0715733ca6f，唯一标示
			cmdAdd netNS: /var/run/netns/testing		# 设置网络命名空间的路径与名称 ex: ip netns add testing or ip netns delete testing
		    cmdAdd ifName: eth0，在容器内创建的接口名称
			cmdAdd args: 使用cni插件 传入环境参数，以分号分隔的字母数字键值对；例如，“FOO=BAR;ABC=123”
			cmdAdd path: ./bin， cni插件 可执行文件路径
			cmdAdd stdin: {\"bridge\":\"jtthink0\",\"cniVersion\":\"0.4.0\",\"ipam\":{\"dataDir\":\"/tmp/cni-host\",\"routes\":[{\"dst\":\"0.0.0.0/0\"}],\"subnet\":\"10.16.0.0/16\",\"type\":\"host-local\"},\"name\":\"mynet\",\"type\":\"jtthink\"}
	*/

	// 获取CNI conf 配置文件对象
	cfg, err := config.LoadCNIConfig(args.StdinData)
	if err != nil {
		fmt.Println("config error: ", err)
		return err
	}

	result := &current.Result{
		CNIVersion: current.ImplementedSpecVersion,
	}

	cniArgs := veth.ParseArgs(args.Args)

	// 用ipam分配ip地址
	if cfg.IPAM.Type != "" {
		r, err := veth.Ipam(cfg)
		if err != nil {
			fmt.Println("ipam error: ", err)
			return err
		}
		ipamRet, err := current.NewResultFromResult(r)
		if err != nil {
			fmt.Println("ipamRet error: ", err)
			return err
		}
		// 返回的结果赋值
		result.IPs = ipamRet.IPs
		result.DNS = ipamRet.DNS
		result.Routes = ipamRet.Routes
	}

	// 创建或更新网桥
	var br *netlink.Bridge
	if br, err = bridge.CreateOrUpdateBridge(cfg.Bridge); err != nil {
		fmt.Println("bridge error: ", err)
		return err
	}

	// 创建veth设备
	err = veth.CreateVeth(args.Netns, result.IPs[0].Address.String(), br, cniArgs.PodName, cniArgs.ContainerID)
	if err != nil {
		fmt.Println("veth error: ", err)
		return err
	}

	return types.PrintResult(result, cfg.CNIVersion)

}

// cmdDel CNI del方法
func cmdDel(args *skel.CmdArgs) error {

	// 获取CNI conf 配置文件对象
	cfg, err := config.LoadCNIConfig(args.StdinData)
	if err != nil {
		fmt.Println("config error: ", err)
		return err
	}

	cniArgs := veth.ParseArgs(args.Args)

	// 释放ipam ip，不需要手动删除ipam文件夹
	err = veth.ReleaseIP(cfg)
	if err != nil {
		fmt.Println("release ip err: ", err)
		return err
	}

	// TODO 删除veth pair
	err = veth.DelVeth(cniArgs.PodName)
	if err != nil {
		return err
	}

	return nil

}

func cmdCheck(*skel.CmdArgs) error {
	return nil
}
