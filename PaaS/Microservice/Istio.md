# Istio #

## 架构 ##

![](img/Istio_arch.svg)

Istio 服务网格逻辑上分为数据平面和控制平面。

- 数据平面由一组以 sidecar 方式部署的智能代理（Envoy）组成。这些代理可以调节和控制微服务及 Mixer 之间所有的网络通信。
 
- 控制平面负责管理和配置代理来路由流量。此外控制平面配置 Mixer 以实施策略和收集遥测数据。

**组件功能**

- Envoy

	参考本目录Envoy内容

- Mixer

	Mixer是一个独立于平台的组件，负责在服务网格上执行访问控制和使用策略，并从Envoy代理和其他服务收集遥测数据。代理提取请求级属性，发送到Mixer进行评估。Mixer中包括一个灵活的插件模型，使其能够接入到各种主机环境和基础设施后端，从这些细节中抽象出Envoy代理和Istio管理的服务。

- Pilot

	Pilot为Envoy sidecar提供服务发现功能，为智能路由（例如 A/B 测试、金丝雀部署等）和弹性（超时、重试、熔断器等）提供流量管理功能。它将控制流量行为的高级路由规则转换为特定于 Envoy 的配置，并在运行时将它们传播到sidecar。Pilot将平台特定的服务发现机制抽象化并将其合成为符合Envoy数据平面API的任何sidecar都可以使用的标准格式。这种松散耦合使得Istio能够在多种环境下运行（例如，Kubernetes、Consul、Nomad），同时保持用于流量管理的相同操作界面。

- Citadel

	Citadel通过内置身份和凭证管理赋能强大的服务间和最终用户身份验证。可用于升级服务网格中未加密的流量，并为运维人员提供基于服务标识而不是网络控制的强制执行策略的能力。从0.5版本开始，Istio支持基于角色的访问控制，以控制谁可以访问您的服务，而不是基于不稳定的三层或四层网络标识。

- Galley
- 
	Galley代表其他的Istio控制平面组件，用来验证用户编写的Istio API配置。随着时间的推移，Galley将接管Istio获取配置、 处理和分配组件的顶级责任。它将负责将其他的Istio组件与从底层平台（例如 Kubernetes）获取用户配置的细节中隔离开来。


*istio-init容器*

该容器就是通过修改iptables规则让Envoy代理可以拦截所有的进出Pod的流量，即将入站流量重定向到Sidecar，再拦截应用容器的出站流量经过Sidecar处理后再出站。该容器的入口脚本为[istio-iptables.sh](https://github.com/istio/istio/blob/master/tools/packaging/common/istio-iptables.sh)，脚本的使用方法(部分，摘自[https://jimmysong.io/posts/envoy-sidecar-injection-in-istio-service-mesh-deep-dive/](https://jimmysong.io/posts/envoy-sidecar-injection-in-istio-service-mesh-deep-dive/),完整使用方法可阅读脚本)：

```
$ istio-iptables.sh -p PORT -u UID -g GID [-m mode] [-b ports] [-d ports] [-i CIDR] [-x CIDR] [-h]
  -p: 指定重定向所有 TCP 流量的 Envoy 端口（默认为 $ENVOY_PORT = 15001）
  -u: 指定未应用重定向的用户的 UID。通常，这是代理容器的 UID（默认为 $ENVOY_USER 的 uid，istio_proxy 的 uid 或 1337）
  -g: 指定未应用重定向的用户的 GID。（与 -u param 相同的默认值）
  -m: 指定入站连接重定向到 Envoy 的模式，“REDIRECT” 或 “TPROXY”（默认为 $ISTIO_INBOUND_INTERCEPTION_MODE)
  -b: 逗号分隔的入站端口列表，其流量将重定向到 Envoy（可选）。使用通配符 “*” 表示重定向所有端口。为空时表示禁用所有入站重定向（默认为 $ISTIO_INBOUND_PORTS）
  -d: 指定要从重定向到 Envoy 中排除（可选）的入站端口列表，以逗号格式分隔。使用通配符“*” 表示重定向所有入站流量（默认为 $ISTIO_LOCAL_EXCLUDE_PORTS）
  -i: 指定重定向到 Envoy（可选）的 IP 地址范围，以逗号分隔的 CIDR 格式列表。使用通配符 “*” 表示重定向所有出站流量。空列表将禁用所有出站重定向（默认为 $ISTIO_SERVICE_CIDR）
  -x: 指定将从重定向中排除的 IP 地址范围，以逗号分隔的 CIDR 格式列表。使用通配符 “*” 表示重定向所有出站流量（默认为 $ISTIO_SERVICE_EXCLUDE_CIDR）。

环境变量位于 $ISTIO_SIDECAR_CONFIG（默认在：/var/lib/istio/envoy/sidecar.env）
```