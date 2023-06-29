# cni-interface-study

### CNI 配置文件
[CNI 配置文件参考](./conf/10-mynet.conf)
```json
{
  "cniVersion": "0.4.0",  # cni版本
  "name": "mynet",        # cni插件名
  "type": "mycniplugin",  # cni二进制文件
  "bridge":"mycni0",      # 网桥名,
  # 调用ipam插件
  "ipam": {
        "type": "host-local",
        "subnet": "10.16.0.0/16",
        "dataDir": "/tmp/cni-host",
        "routes": [
            { "dst": "0.0.0.0/0" }
        ]
  }
}
```

#### 项目测试
- 1. 创建 netns
```bash
ip netns add testing
```    
- 2. 执行 cnitool 测试插件(需要先安装 cnitool 工具)
```bash
NETCONFPATH=./conf  CNI_PATH=./bin cnitool add mynet /var/run/netns/testing
[root@VM-0-8-centos cni_practice]# NETCONFPATH=./conf  CNI_PATH=./bin cnitool add mynet /var/run/netns/testing
{
    "cniVersion": "0.4.0",
    "ips": [
        {
            "version": "4",
            "address": "10.16.0.4/16",
            "gateway": "10.16.0.1"
        }
    ],
    "routes": [
        {
            "dst": "0.0.0.0/0"
        }
    ],
    "dns": {}
}

``` 
- 3. 查看网桥与 veth 设备
```bash
[root@VM-0-8-centos cni_practice]# ifconfig
mycni0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.16.0.1  netmask 255.255.0.0  broadcast 10.16.255.255
        inet6 fe80::4c6b:8bff:feb2:e6cd  prefixlen 64  scopeid 0x20<link>
        ether fe:71:f7:48:52:63  txqueuelen 0  (Ethernet)
        RX packets 21  bytes 1464 (1.4 KiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 8  bytes 656 (656.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
pod1234: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet6 fe80::fc71:f7ff:fe48:5263  prefixlen 64  scopeid 0x20<link>
        ether fe:71:f7:48:52:63  txqueuelen 0  (Ethernet)
        RX packets 7  bytes 586 (586.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 7  bytes 586 (586.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0        
```    
- 4. 测试
```bash 
# ping netns内的网络是通的
[root@VM-0-8-centos cni_practice]# ping 10.16.0.4
PING 10.16.0.4 (10.16.0.4) 56(84) bytes of data.
64 bytes from 10.16.0.4: icmp_seq=1 ttl=64 time=0.060 ms
64 bytes from 10.16.0.4: icmp_seq=2 ttl=64 time=0.049 ms
# ping 网桥是通的
[root@VM-0-8-centos cni_practice]# ping 10.16.0.1
PING 10.16.0.1 (10.16.0.1) 56(84) bytes of data.
64 bytes from 10.16.0.1: icmp_seq=1 ttl=64 time=0.035 ms
64 bytes from 10.16.0.1: icmp_seq=2 ttl=64 time=0.038 ms
# ping 宿主机网桥是通的
[root@VM-0-8-centos cni_practice]# ip netns exec testing ping 10.16.0.1
PING 10.16.0.1 (10.16.0.1) 56(84) bytes of data.
64 bytes from 10.16.0.1: icmp_seq=1 ttl=64 time=0.046 ms
64 bytes from 10.16.0.1: icmp_seq=2 ttl=64 time=0.051 ms
# 由 netns 中可 ping 通宿主机网络
[root@VM-0-8-centos cni_practice]# ip netns exec testing ping 10.0.0.8
PING 10.0.0.8 (10.0.0.8) 56(84) bytes of data.
64 bytes from 10.0.0.8: icmp_seq=1 ttl=64 time=0.048 ms
64 bytes from 10.0.0.8: icmp_seq=2 ttl=64 time=0.051 ms
```    
- 5. 删除设备
```bash
NETCONFPATH=./conf  CNI_PATH=./bin cnitool del mynet /var/run/netns/testing
```    

#### 项目测试用
```bash
# 删除遗留
ip netns delete testing && rm /tmp/cni-host/mynet -fr

go build -o ./bin/mycniplugin ./cmd/cni/main.go
# cnitool 测试
NETCONFPATH=./conf  CNI_PATH=./bin cnitool add mynet /var/run/netns/testing
NETCONFPATH=./conf  CNI_PATH=./bin cnitool del mynet /var/run/netns/testing
NETCONFPATH=./conf  CNI_PATH=./bin cnitool check mynet /var/run/netns/testing
```