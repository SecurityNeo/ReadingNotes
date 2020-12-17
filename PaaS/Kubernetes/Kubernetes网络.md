# Kubernetes网络 #

## 基础知识 ##

### veth-pair ###

veth-pair 是成对出现的一种虚拟网络设备，一端连接着协议栈，一端彼此相连，数据从一端出，从另一端进。

1. 创建一对名为 veth0 和 veth1 的 veth 接口：

	`ip link add veth0 type veth peer name veth1`

2. 创建 ns1 网络命名空间：

	`ip netns add ns1`

3. 将 veth0 接口加到 ns1 网络命名空间里：

	`ip link set veth0 netns ns1`
 
4. 为 veth0 接口配置 IP 地址：
 
	`ip -n ns1 addr add 10.1.1.1/24 dev veth0`
 
5. 将 veth0 接口和 lo 口 up 起来：

	`ip -n ns1 link set veth0 up`
	`ip -n ns1 link set lo up`

```shell
[root@VM-0-4-centos ~]# ip -n ns1 addr show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
5: veth0@if4: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state LOWERLAYERDOWN group default qlen 1000
    link/ether 52:a3:29:82:70:75 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.1.1.1/24 scope global veth0
       valid_lft forever preferred_lft forever
[root@VM-0-4-centos ~]# ip netns exec ns1 ping -c2 10.1.1.1
PING 10.1.1.1 (10.1.1.1) 56(84) bytes of data.
64 bytes from 10.1.1.1: icmp_seq=1 ttl=64 time=0.041 ms
64 bytes from 10.1.1.1: icmp_seq=2 ttl=64 time=0.059 ms

```

现在配置另一个接口veth1

1. 创建 ns2 网络命名空间：

	`ip netns add ns2`

2. 将 veth1 接口加到 ns2 网络命名空间里：

	`ip link set veth1 netns ns2`

3. 为 veth1 接口配置 IP 地址：

	`ip -n ns2 addr add 10.2.1.1/24 dev veth1`

4. 将 veth1 接口和 lo 口 up 起来：

	`ip -n ns2 link set veth1 up`
	`ip -n ns2 link set lo up`

```shell
[root@VM-0-4-centos ~]# ip -n ns2 addr show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
4: veth1@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether d6:c4:03:51:f9:98 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.2.1.1/24 scope global veth1
       valid_lft forever preferred_lft forever
    inet6 fe80::d4c4:3ff:fe51:f998/64 scope link
       valid_lft forever preferred_lft forever
[root@VM-0-4-centos ~]# ip netns exec ns2 ping -c2 10.2.1.1
PING 10.2.1.1 (10.2.1.1) 56(84) bytes of data.
64 bytes from 10.2.1.1: icmp_seq=1 ttl=64 time=0.040 ms
64 bytes from 10.2.1.1: icmp_seq=2 ttl=64 time=0.060 ms
```

```shell
[root@VM-0-4-centos ~]# ip netns exec ns1 ethtool -S veth0
NIC statistics:
     peer_ifindex: 4
[root@VM-0-4-centos ~]# ip netns exec ns2 ethtool -S veth1
NIC statistics:
     peer_ifindex: 5
```

此时两个命名空间的路由如下：

```shell
[root@VM-0-4-centos ~]# ip netns exec ns1 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth0
[root@VM-0-4-centos ~]# ip netns exec ns2 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.2.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth1

```

此时两个接口之间没法通信，我们需要分别为其添加路由

```shell
[root@VM-0-4-centos ~]# ip -n ns1 route add 10.2.1.0/24 dev veth0
[root@VM-0-4-centos ~]# ip -n ns2 route add 10.1.1.0/24 dev veth1
[root@VM-0-4-centos ~]# ip netns exec ns1 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth0
10.2.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth0
[root@VM-0-4-centos ~]# ip netns exec ns2 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth1
10.2.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth1
```

```shell
[root@VM-0-4-centos ~]# ip netns exec ns1 ping -c2 10.2.1.1
PING 10.2.1.1 (10.2.1.1) 56(84) bytes of data.
64 bytes from 10.2.1.1: icmp_seq=1 ttl=64 time=0.068 ms
64 bytes from 10.2.1.1: icmp_seq=2 ttl=64 time=0.074 ms

--- 10.2.1.1 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 999ms
rtt min/avg/max/mdev = 0.068/0.071/0.074/0.003 ms
[root@VM-0-4-centos ~]# ip netns exec ns2 ping -c2 10.1.1.1
PING 10.1.1.1 (10.1.1.1) 56(84) bytes of data.
64 bytes from 10.1.1.1: icmp_seq=1 ttl=64 time=0.052 ms
64 bytes from 10.1.1.1: icmp_seq=2 ttl=64 time=0.071 ms

--- 10.1.1.1 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 999ms
rtt min/avg/max/mdev = 0.052/0.061/0.071/0.012 ms

```

### ipset ###

### conntract ###

### iptables ###

### ipvs ###