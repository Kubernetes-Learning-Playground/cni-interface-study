package veth

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"net"
	"os"

	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// CreateVeth 创建veth
func CreateVeth(nspath string, addrstr string, br *netlink.Bridge, vethHost, vethContainer string) error {

	if oldHostVeth, err := netlink.LinkByName(vethHost); err == nil {
		if err = netlink.LinkDel(oldHostVeth); err != nil {
			return errors.Wrapf(err, "failed to delete old hostVeth %v", err)
		}
	}

	// TODO: 名字需要修改
	//var vethHost, vethContainer = RandomVethName(), RandomVethName()
	vethpeer := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: vethHost, MTU: 1500},
		PeerName:  vethContainer, //随机名称  一般是veth@xxxx  建立在宿主机上
	}

	// 执行ip link add
	err := netlink.LinkAdd(vethpeer)
	if err != nil {
		return errors.Wrapf(err, "netlink.LinkAdd error")
	}

	// 宿主机

	// 宿主机veth端
	myveth_host, err := netlink.LinkByName(vethHost)
	if err != nil {
		return errors.Wrapf(err, "netlink.LinkByName error")
	}
	// 挂bridge
	err = netlink.LinkSetMaster(myveth_host, br)
	if err != nil {
		return errors.Wrapf(err, "netlink.LinkSetMaster error")
	}
	// 启动veth
	err = netlink.LinkSetUp(myveth_host)
	if err != nil {
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}

	ns, err := netns.GetFromPath(nspath)
	if err != nil {
		fmt.Println("netns.GetFromPath err: ", err)
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}
	defer ns.Close()

	//获取 容器里面的 veth设备
	myvethContainer, err := netlink.LinkByName(vethContainer)
	if err != nil {
		fmt.Println("netlink.LinkByName err: ", err)
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}

	// 容器中

	// ip link set xxx netns nsname
	// 把另一端veth放入容器ns
	err = netlink.LinkSetNsFd(myvethContainer, int(ns))
	if err != nil {
		fmt.Println("netlink.LinkSetNsFd err: ", err)
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}

	// 进入ns空间
	err = netns.Set(ns)
	if err != nil {
		fmt.Println("netns.Set add err: ", err)
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}
	// 获取ns的veth设备
	myvethContainer, err = netlink.LinkByName(vethContainer)
	if err != nil {
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}

	// 设置地址
	addr, _ := netlink.ParseAddr(addrstr)
	//设置IP地址
	err = netlink.AddrAdd(myvethContainer, addr)
	if err != nil {
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}
	// ns中veth设备名称改为eth0
	_ = netlink.LinkSetName(myvethContainer, "eth0")
	// 启动
	err = netlink.LinkSetUp(myvethContainer)
	if err != nil {
		return errors.Wrapf(err, "netlink.LinkSetUp error")
	}
	return addRoute()
}

// addRoute 为ns内部的网络添加路由，才能让容器内的ns与容器外的互通
// ex: ip netns exec testing ping 10.0.0.8
// ex: ping 10.16.0.2
// ip netns exec testing route -n 查看ns的路由
func addRoute() error {
	route := &netlink.Route{
		Dst: &net.IPNet{
			IP:   net.IPv4(0, 0, 0, 0),
			Mask: net.IPv4Mask(0, 0, 0, 0),
		},
		Gw: net.IPv4(10, 16, 0, 1), //网关地址  -- 网桥IP
	}
	return netlink.RouteAdd(route)
}

// DelVeth 删除veth设备
func DelVeth(hostVethName string) error {
	// 删除veth pair
	if err := ip.DelLinkByName(hostVethName); err != nil {
		klog.Error("del link err: ", err)
		return errors.Wrapf(err, "ip.DelLinkByName error")
	}

	// TODO 不需要进入ns删除，
	// 进入namespace 删除设备
	//_, err := ns.GetNS(nspath)
	//if err != nil {
	//	return errors.Wrapf(err, "ns.GetNS error")
	//}
	//defer netns.Close()

	//err = delVeth(netNs, "eth0")
	//if err != nil {
	//	return errors.Wrapf(err, "delVeth error")
	//}

	return nil
}

func delVeth(netns ns.NetNS, ifName string) error {
	return netns.Do(func(ns.NetNS) error {
		l, err := netlink.LinkByName(ifName)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		return netlink.LinkDel(l)
	})
}
