# cni-interface-study

项目测试用
```bash
ip link del pod1234
ip netns add testing
ip netns delete testing
rm  /tmp/cni-host/mynet -fr
NETCONFPATH=./conf  CNI_PATH=./bin cnitool add mynet /var/run/netns/testing
go build -o ./bin/jtthink ./cmd/cni/main.go
NETCONFPATH=./conf  CNI_PATH=./bin cnitool del mynet /var/run/netns/testing
```