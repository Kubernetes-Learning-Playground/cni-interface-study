package bridge

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

// CreateOrUpdateBridge 创建或更新网桥
func CreateOrUpdateBridge(brName, subnet string) (*netlink.Bridge, error) {
	// 按照名称取出设备
	link, err := netlink.LinkByName(brName)
	if err != nil {
		// 如果是找不到的错误，则创建新的设备
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			br := &netlink.Bridge{
				LinkAttrs: netlink.LinkAttrs{Name: brName, MTU: 1500},
			}

			// 加入网桥对象
			if err := netlink.LinkAdd(br); err != nil {
				return nil, err
			}

			var addr *netlink.Addr
			if addr, err = netlink.ParseAddr(subnet); err != nil {
				fmt.Println("ParseAddr add err: ", err)
				return nil, err
			}

			if err = netlink.AddrAdd(br, addr); err != nil {
				fmt.Println("addr add err: ", err)
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
		return nil, fmt.Errorf("error bridge")
	}
}
