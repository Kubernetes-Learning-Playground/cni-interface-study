package veth

import (
	"crypto/rand"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// CreateVeth 创建veth
func CreateVeth(nspath string, addrstr string, br *netlink.Bridge) error {
	// 生成veth设备对名称
	// TODO: 名字需要修改
	var veth_host, veth_container = RandomVethName(), RandomVethName()
	vethpeer := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: veth_host, MTU: 1500},
		PeerName:  veth_container, //随机名称  一般是veth@xxxx  建立在宿主机上
	}

	// 执行ip link add
	err := netlink.LinkAdd(vethpeer)
	if err != nil {
		return err
	}

	// 宿主机

	// 宿主机veth端
	myveth_host, err := netlink.LinkByName(veth_host)
	if err != nil {
		return err
	}
	// 挂bridge
	err = netlink.LinkSetMaster(myveth_host, br)
	if err != nil {
		return err
	}
	// 启动veth
	err = netlink.LinkSetUp(myveth_host)
	if err != nil {
		return err
	}

	ns, err := netns.GetFromPath(nspath)
	if err != nil {
		return err
	}
	defer ns.Close()

	//获取 容器里面的 veth设备
	myveth_container, err := netlink.LinkByName(veth_container)
	if err != nil {
		return err
	}

	// 容器中

	// ip link set xxx netns nsname
	// 把另一端veth放入容器ns
	err = netlink.LinkSetNsFd(myveth_container, int(ns))
	if err != nil {
		return err
	}

	// 进入ns空间
	err = netns.Set(ns)
	if err != nil {
		return err
	}
	// 获取ns的veth设备
	myveth_container, err = netlink.LinkByName(veth_container)
	if err != nil {
		return err
	}

	// 设置地址
	addr, _ := netlink.ParseAddr(addrstr)
	//设置IP地址
	err = netlink.AddrAdd(myveth_container, addr)
	if err != nil {
		return err
	}
	// ns中veth设备名称改为eth0
	_ = netlink.LinkSetName(myveth_container, "eth0")
	// 启动
	err = netlink.LinkSetUp(myveth_container)
	if err != nil {
		return err
	}
	return addRoute()
}

// RandomVethName 生成VethName
func RandomVethName() string {
	entropy := make([]byte, 4)
	rand.Read(entropy)
	return fmt.Sprintf("jtveth%x", entropy)
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
