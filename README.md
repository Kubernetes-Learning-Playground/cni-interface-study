# cni-interface-study

项目测试用
```bash
ip netns add testing
ip netns delete testing
rm  /tmp/cni-host/mynet -fr
NETCONFPATH=./conf  CNI_PATH=./bin cnitool add mynet /var/run/netns/testing
```