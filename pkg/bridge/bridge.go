package bridge

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

// CreateOrUpdateBridge 创建或更新网桥
func CreateOrUpdateBridge(br string) (*netlink.Bridge, error) {
	// 根据名称取出设备
	link, err := netlink.LinkByName(br)
	if err != nil {
		// 如果是找不到的错误，则创建新的设备
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			br := &netlink.Bridge{
				LinkAttrs: netlink.LinkAttrs{Name: br, MTU: 1500},
			}
			// 加入网桥对象
			if err := netlink.LinkAdd(br); err != nil {
				return nil, err
			}
			// 暂时写死addr
			var addr *netlink.Addr
			if addr, err = netlink.ParseAddr("10.16.0.1/16"); err != nil {
				return nil, err
			}
			if err = netlink.AddrAdd(br, addr); err != nil {
				return nil, err
			}
			if err = netlink.LinkSetUp(br); err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	if br, ok := link.(*netlink.Bridge); ok {
		return br, nil
	} else {
		return nil, fmt.Errorf("错误的网桥对象")
	}

}
