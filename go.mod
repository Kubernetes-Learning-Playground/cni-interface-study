module github/mycni/cni_practice

go 1.18

require (
	github.com/containernetworking/cni v1.1.2
	github.com/containernetworking/plugins v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/vishvananda/netlink v1.2.1-beta.2
	github.com/vishvananda/netns v0.0.4
	k8s.io/klog/v2 v2.100.1
)

require (
	github.com/coreos/go-iptables v0.6.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/safchain/ethtool v0.2.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
)
