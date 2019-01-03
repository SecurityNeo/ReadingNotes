# VRRP #

[RFC-2338](https://tools.ietf.org/html/rfc2338)

[RFC-3768](https://tools.ietf.org/html/rfc3768)


## 基本概念 ##

**VRRP路由器(VRRP Router)**

运行VRRP协议一个或多个实例的路由器设备

**虚拟路由器(Virtual Router)**

由一个Master路由器和多个Backup路由器组成，也被称为VRRP备份组。也是一个共享局域网内主机的默认网关

**Master路由器**

转发报文或者应答ARP请求的VRRP路由器，承担着流量的转发

**Backup路由器**

一组没有承担转发任务的VRRP路由器，当一个虚拟路由器中的Master路由器出现故障时，它们可以通过竞选成为新的Master

**虚拟IP地址(Virtual IP Address)**

虚拟路由器的IP地址，一个虚拟路由器可以拥有一个或多个虚拟IP地址

**主IP地址(Primary IP Address)**

从接口的真实IP地址中选出来的一个主用IP地址，通常选择配置的第一个IP地址。 VRRP通告报文的源地址总是主IP地址

**虚拟MAC地址**

虚拟路由器根据虚拟路由器ID生成的MAC地址，当虚拟路由器回应ARP请求时，使用虚拟MAC地址，而不是接口的真实MAC地址回应ARP请求。虚拟MAC地址组成方式是：00-00-5E-00-01-{VRID}，前三个字节00-00-5E是IANA组织分配的，接下来的两个字节00-01是为VRRP协议指定的，最后的VRID是虚拟路由器标识，取值范围[1，255] 

**VRID**

虚拟路由器标识，在同一个VRRP组内的路由器必须有相同的VRID，通过VRID表明自己属于哪个VRRP组

## 通告报文结构 ##

### VRRP通告报文结构 ###

VRRP协议只有一种报文，即VRRP报文。VRRP报文用来将Master设备的优先级和状态通告给同一虚拟路由器的所有VRRP路由器。VRRP报文封装在IP报文中，报文结构如下：

![报文结构](img/VRRP.png)

**VRRP报文字段简介**：

- Version

协议版本号，4位，在RFC3768中定义为2

- Type

报文类型，4位，只会为1，表示Advertisement


- Virtual Rtr ID

虚拟路由器ID，8位，取值范围是1～255


- Priority

发送报文的VRRP路由器在虚拟路由器中的优先级，8位。取值范围是0～255，其中可用的范围是1～254。0表示设备停止参与VRRP，用来使备份路由器尽快成为主路由器，而不必等到计时器超时；255则保留给IP地址拥有者。缺省值是100。


- Count IP Addrs

VRRP中IP地址个数，8位。


- Authentication Type

验证类型，8位，RFC2338定义的取值为：

0 - No Authentication

1 - Simple Text Password

2 - IP Authentication Header

随后的RFC3768中将Authentication Type取值变更为（即取消认证，因为这些认证方式并不能提供真正的安全）：

0 - No Authentication

1 - Reserved

2 - Reserved


- Adver Int

通告包的发送间隔时间，8位，单位是秒，默认为1秒


- Checksum

校验和，16位，只针对VRRP数据部分进行校验，不包括IP头部

- IP Address(es)

虚拟路由器IP地址，数量由Count IP Addrs决定，这些信息主要用于发现并修复误配置路由器


- Authentication Data

RFC3768中定义该字段只是为了和老版本兼容，必须置0

**VRRP报文中IP字段**


- 源IP地址

Master路由器发送包的物理接口IP地址


- 目的IP地址

IP组播地址224.0.0.18


- TTL

必须为255


- 协议号

0x70（十进制为112）

## VRRP协议状态机 ##

VRRP协议中定义了三种状态机：初始状态（Initialize）、活动状态（Master）、备份状态（Backup）。其中，只有处于活动状态的设备才可以转发那些发送到虚拟IP地址的报文。



