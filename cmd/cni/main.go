package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/040"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
	"github.com/vishvananda/netlink"
	"github/mycni/cni_practice/pkg/bridge"
	"github/mycni/cni_practice/pkg/config"
	"github/mycni/cni_practice/pkg/veth"
	"k8s.io/klog/v2"
)

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("mycniplugin"))
}

// cmdAdd CNI add方法
func cmdAdd(args *skel.CmdArgs) error {
	// 使用 cnitool 打印 log 会报错。
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
		cmdAdd args: 使用cni插件 传入环境参数，以分号分隔的字母数字键值对；例如，“K8S_POD_NAMESPACE=default;K8S_POD_NAME=pod-test”
		cmdAdd path: ./bin， cni插件 可执行文件路径
		cmdAdd stdin: {\"bridge\":\"jtthink0\",\"cniVersion\":\"0.4.0\",\"ipam\":{\"dataDir\":\"/tmp/cni-host\",\"routes\":[{\"dst\":\"0.0.0.0/0\"}],\"subnet\":\"10.16.0.0/16\",\"type\":\"host-local\"},\"name\":\"mynet\",\"type\":\"jtthink\"}
	*/

	// 获取 CNI conf 配置文件对象
	cfg, err := config.LoadCNIConfig(args.StdinData)
	if err != nil {
		klog.Error("config error: ", err)
		return err
	}

	result := &current.Result{
		CNIVersion: current.ImplementedSpecVersion,
	}

	// 解析参数
	cniArgs := veth.ParseArgs(args.Args)

	// 用 ipam 分配 ip 地址
	if cfg.IPAM.Type != "" {
		r, err := veth.Ipam(cfg)
		if err != nil {
			klog.Error("ipam error: ", err)
			return err
		}
		ipamRet, err := current.NewResultFromResult(r)
		if err != nil {
			klog.Error("ipamRet error: ", err)
			return err
		}
		// 返回的结果赋值
		result.IPs = ipamRet.IPs
		result.DNS = ipamRet.DNS
		result.Routes = ipamRet.Routes
	}

	// 创建或更新网桥
	var br *netlink.Bridge
	if br, err = bridge.CreateOrUpdateBridge(cfg.Bridge, cfg.IPAM.Subnet); err != nil {
		klog.Error("bridge error: ", err)
		return err
	}

	// 创建 veth 设备
	err = veth.CreateVeth(args.Netns, result.IPs[0].Address.String(), br, cniArgs.PodName, cniArgs.ContainerID)
	if err != nil {
		klog.Error("veth error: ", err)
		return err
	}

	return types.PrintResult(result, cfg.CNIVersion)

}

// cmdDel CNI del方法
func cmdDel(args *skel.CmdArgs) error {

	// 获取 CNI conf 配置文件
	cfg, err := config.LoadCNIConfig(args.StdinData)
	if err != nil {
		klog.Error("config error: ", err)
		return err
	}

	cniArgs := veth.ParseArgs(args.Args)

	// 释放 ipam ip，不需要手动删除 ipam 文件夹
	err = veth.ReleaseIP(cfg)
	if err != nil {
		klog.Error("release ip err: ", err)
		return err
	}

	// TODO 删除 veth pair
	err = veth.DelVeth(cniArgs.PodName)
	if err != nil {
		return err
	}

	return nil

}

func cmdCheck(*skel.CmdArgs) error {
	return nil
}
