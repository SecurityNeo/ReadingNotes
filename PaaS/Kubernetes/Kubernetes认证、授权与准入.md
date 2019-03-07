# Kubernetes认证、授权和准入 #

## 认证 ##

kubernetes认证有很多种：CA证书认证、Token认证、Basic认证、Webhook Token认证、Bootstrap Token认证（也是一种Token认证）、OpenID Connect Tokens认证等。可以同时配置多种认证方式，只要其中任意一个方式认证通过即可。

- CA证书认证

	在apiserver的启动参数中指定`--client_ca_file=SOMEFILE`即可开启CA证书认证，被引用的文件中包含验证client的证书，如果被验证通过，那么这个验证记录中的主体对象将会作为请求的username。

- Token认证

	在apiserver的启动参数中指定`--token_auth_file=SOMEFILE`即可开启静态Token认证。 token文件包含三列：token，username，userid，第四列为可选group名，多个group需用逗号分隔，并且需用双引号，如`token,user,uid,"group1,group2,group3"`。当使用token作为验证方式时，在对apiserver的http请求中，增加一个Header字段：Authorization ，将它的值设置为`Bearer SOMETOKEN`。

- Basic认证

	在apiserver的启动参数中指定`--basic_auth_file=SOMEFILE`即可开启Basic认证，如果更改了文件中的密码，只有重新启动apiserver才能使其重新生效。其文件的基本格式包含三列：password，username，userid,在Kubernetes 1.6+ 版本可以指定一个可选的第四列，使用逗号来分隔group名，如果有多个组，则该组必须使用双引号。当使用此作为认证方式时，在对apiserver的http请求中，需增加一个头部：Authorization，值为：`Basic BASE64ENCODED(USER:PASSWORD)`。

- Bootstrap Token认证

	Kubernetes包括一个dynamically-managed的Bearer token类型，称为Bootstrap Token。这些tokens作为Secret存储在kube-system namespace中，可以动态管理和创建。Controller Manager包含一个TokenCleaner Controller，如果到期时可以删除Bootstrap Tokens Controller。Tokens的形式是[a-z0-9]{6}.[a-z0-9]{16}。第一个组件是Token ID，第二个组件是Token Secret。可以在HTTP header中指定Token:`Authorization: Bearer 781292.db7bc3a58fc5f07e`。Bootstrap Token当前仅用作kubelet的客户端证书。

	kubelet发起的CSR请求都是由 controller manager 来做实际签署的，对于 controller manager 来说，TLS bootstrapping 下 kubelet 发起的 CSR 请求大致分为以下三种

		- nodeclient
		
			kubelet以O=system:nodes和CN=system:node:(node name) 形式发起的CSR请求，仅在第一次启动时会产生

		- selfnodeclient
		
			kubelet client renew自己的证书发起的CSR请求(与上一个证书就有相同的O和CN)，kubelet renew自己作为client apiserver通讯时使用的证书产生

		- selfnodeserver
		
			kubelet server renew自己的证书发起的CSR请求，kubelet首次申请或后续renew自己的10250 api端口证书时产生

- ServiceAccount Token认证

	ServiceAccount是一个自动启用的认证器，它使用被签名的Bearer Token对请求进行认证，该插件有两个可选参数：

	`--service-account-key-file`: 一个包含签名bearer token的PEM编码文件。如果未指定，将使用API server的TLS私钥。
	
	`--service-account-lookup`: 如果启用，从API中删除掉的token将被撤销。

	ServiceAccount是限定命名空间的，在创建命名空间的时候，kubernetes会为每一个命名空间创建一个默认的ServiceAccount，这个默认的ServiceAccount只能访问该命名空间内的资源。Service Account和Pod、Service、Deployment一样是 Kubernetes 群中的一种资源，我们也可以手动创建Service Account。Service Account使用用户名进行验证`system:serviceaccount:(NAMESPACE):(SERVICEACCOUNT)`，并分配给组`system:serviceaccounts`和`system:serviceaccounts:(NAMESPACE)`。当令牌通过认证后，请求的用户名将被设置为 `system:serviceaccount:(NAMESPACE):(SERVICEACCOUNT)`，请求的组名为：`system:serviceaccounts`和`system:serviceaccounts:(NAMESPACE)`。

	ServiceAccount主要包含了三个内容：命名空间、令牌 和 CA。它们都通过挂接的方式保存在pod的文件系统中，他们的路径是：

	- 令牌: `/var/run/secrets/kubernetes.io/serviceaccount/token`，这是apiserver使用base64编码通过私钥签发的令牌
	- CA： `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`
	- 命名空间： `/var/run/secrets/kubernetes.io/serviceaccount/namespace`,也使用base64编码。


- Webhook Token认证

	Webhook Token认证方式让用户使用自己的认证方式，用户需按照约定的请求格式和应答格式提供HTTPS服务，当用户把Bearer Token放到请求的头部，kubernetes会把token发送给事先配置的地址进行认证，如果认证结果成功，则认为请求用户合法。 相关配置参数：

	- authentication-token-webhook-config-file ：配置文件描述了如何访问远程webhook服务

	- authentication-token-webhook-cache-ttl ：缓存认证时间，默认是两分钟 


- OpenID Connect Tokens（OIDC）

	OpenID Connect是由一些云提供商支持的OAuth2认证机制，常见的有Azure Active Directory，Salesforce和Google。OIDC对OAuth2 协议的主要扩展在于认证成功后，Auth Server除了能授予Access Token之外，还能同时授予ID Token。OIDC ID Token是一种JSON Web Token (JWT) ，这种token中包含一些预定义的域，例如：用户ID，用户分组、email等等… 这些信息非常重要，因为Kubernetes正是把这些信息作为User Account Profile，来进行后续的授权（Authorization）操作的。工作流程如下：

	![](img/oidc.png)

	- 登录identity provider
	- identity provider提供一个access_token，id_token和refresh_token
	- 在使用kubectl时，使用id_token的—token flag或者直接添加到kubeconfig
	- kubectl 在header通过Authorization字段将id_token发送到 API server
	- APIServer检查配置中证书来确保JWT签名的有效性
	- 检查id_token是否过期
	- 确保user被授权
	- 当被授权API Server会返回对kubectl响应
	- kubectl向user反馈

	![](img/OIDC_Parameters.png)

- Keystone认证

	Kubernetes也可以使用Openstack的Keystone组件进行身份认证，涉及两个APIServer参数：`–experimental-keystone-url=<AuthURL>`，`–experimental-keystone-ca-file=SOMEFILE`(开启https时配置)。
	
	参考示例[Kubernetes and Keystone: An integration test passed with flying colors](http://superuser.openstack.org/articles/kubernetes-keystone-integration-test/)

## 授权 ##

Kubernetes中的认证与授权是分开的，授权发生在认证完成之后，认证过程是检验发起API请求的用户是否合法。授权是判断此用户是否有执行该API请求的权限。当配置多个授权模块时，会按顺序检查每个模块，如果有任何模块授权通过，则继续执行下一步的请求。如果所有模块拒绝，则该请求授权失败（返回HTTP 403）。

Kubernetes目前提供以下几种授权模块：

- AlwaysDeny： 表示拒绝所有的请求，仅用作测试。
- AlwaysAllow：允许接收所有请求，如果集群不需要授权流程，则可以采用该策略，这也是Kubernetes的默认配置。
- ABAC： 基于属性的访问控制，表示使用用户配置的授权规则对用户请求进行匹配和控制。
- RBAC：  基于角色的访问控制，允许管理员通过api动态配置授权策略
- Webhook：通过调用外部REST服务对用户进行授权。
- Node：节点授权是一种特殊用途的授权模式，专门授权由kubelet发出的API请求。
- Custom Modules：自定义授权模块。

Kubernetes审查的API请求属性：

- user：身份验证期间提供的user字符串。
- group：经过身份验证的用户所属的组名列表。
- extra：由身份验证层提供的任意字符串键到字符串值的映射。
- API：指示请求是否针对API资源。
- Request path：各种非资源端点的路径，如/api或/healthz。
- API request verb：API动词get，list，create，update，patch，watch，proxy，redirect，delete和deletecollection，用于资源请求。
- HTTP request verb：HTTP动词get，post，put和delete用于非资源请求。
- Resource：正在访问的资源的ID或名称（仅限资源请求） 对于使用get，update，patch和delete动词的资源请求，必须提供资源名称。
- Subresource：正在访问的子资源（仅限资源请求）。
- Namespace：正在访问的对象的名称空间（仅适用于命名空间资源请求）。
- API group：正在访问的API组（仅限资源请求）。空字符串表示核心API组。

