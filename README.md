# cni-interface-study

项目测试用
```bash

ip netns add testing
ip netns delete testing && rm  /tmp/cni-host/mynet -fr

go build -o ./bin/mycniplugin ./cmd/cni/main.go
NETCONFPATH=./conf  CNI_PATH=./bin cnitool add mynet /var/run/netns/testing
NETCONFPATH=./conf  CNI_PATH=./bin cnitool del mynet /var/run/netns/testing
NETCONFPATH=./conf  CNI_PATH=./bin cnitool check mynet /var/run/netns/testing
```