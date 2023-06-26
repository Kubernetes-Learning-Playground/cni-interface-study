package veth

import (
	"crypto/rand"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

/*
 nspath 来自于 args.Netns  网络命名空间的路径 。 譬如 /var/run/netns/testing
addr 在我们分配ip 后 获取
*/
func CreateVeth(nspath string, addrstr string, br *netlink.Bridge) error {
	// 生成veth设备 对  的名称
	var veth_host, veth_container = RandomVethName(), RandomVethName()
	vethpeer := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: veth_host, MTU: 1500},
		PeerName:  veth_container, //随机名称  一般是veth@xxxx  建立在宿主机上
	}

	// 好比执行了ip link add
	err := netlink.LinkAdd(vethpeer)
	if err != nil {
		return err
	}

	// 在宿主机里面的 veth 端
	myveth_host, err := netlink.LinkByName(veth_host)
	if err != nil {
		return err
	}
	err = netlink.LinkSetMaster(myveth_host, br)
	if err != nil {
		return err
	}
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


	// ip link set xxx netns nsname   给 veth设备 设置network namespace
	// 好比 移动到 容器中
	err = netlink.LinkSetNsFd(myveth_container, int(ns))
	if err != nil {
		return err
	}

	err = netns.Set(ns) //设置当前网络命名空间   因为要进入 这个ns 进行操作
	if err != nil {
		return err
	}
	//重新获取 容器内的 veth设备对象
	myveth_container, err = netlink.LinkByName(veth_container)
	if err != nil {
		return err
	}


	addr, _ := netlink.ParseAddr(addrstr)
	//设置IP地址
	err = netlink.AddrAdd(myveth_container, addr)
	if err != nil {
		return err
	}
	// 把容器中 veth设备名称改成 eth0 ,看起来有逼格
	_ = netlink.LinkSetName(myveth_container, "eth0")
	// 启动
	err = netlink.LinkSetUp(myveth_container)
	if err != nil {
		return err
	}
	return addRoute()

}

//随机产生一个VethName
func RandomVethName() string {
	entropy := make([]byte, 4)
	rand.Read(entropy)
	return fmt.Sprintf("jtveth%x", entropy)
}


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