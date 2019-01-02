# VRRP #

[RFC-2338](https://tools.ietf.org/html/rfc2338)

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

