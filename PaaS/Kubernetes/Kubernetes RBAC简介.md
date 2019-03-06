# Kubernetes RBAC简介 #

## APIServer的认证与授权 ##

kubernetes认证主要有五种：CA证书认证、Token认证、Basic认证、Webhook Token认证、Bootstrap Token认证（也是一种Token认证）、OpenID Connect Tokens认证。可以同时配置多种认证方式，只要其中任意一个方式认证通过即可。

- CA证书认证

	在apiserver的启动参数中指定`--client_ca_file=SOMEFILE`即可开启CA证书认证，被引用的文件中包含验证client的证书，如果被验证通过，那么这个验证记录中的主体对象将会作为请求的username。

- Token认证

	在apiserver的启动参数中指定`--token_auth_file=SOMEFILE`即可开启静态Token认证。 token文件包含三列：token，username，userid，第四列为可选group名，多个group需用逗号分隔，并且需用双引号，如`token,user,uid,"group1,group2,group3"`。当使用token作为验证方式时，在对apiserver的http请求中，增加一个Header字段：Authorization ，将它的值设置为`Bearer SOMETOKEN`。

- Basic认证

	在apiserver的启动参数中指定`--basic_auth_file=SOMEFILE`即可开启Basic认证，如果更改了文件中的密码，只有重新启动apiserver才能使其重新生效。其文件的基本格式包含三列：password，username，userid,在Kubernetes 1.6+ 版本可以指定一个可选的第四列，使用逗号来分隔group名，如果有多个组，必须使用双引号。当使用此作为认证方式时，在对apiserver的http请求中，需增加一个头部：Authorization，值为：`Basic BASE64ENCODED(USER:PASSWORD)`。

- Bootstrap Token认证

	Kubernetes包括一个dynamically-managed的Bearer token类型，称为Bootstrap Token。这些tokens作为Secret存储在kube-system namespace中，可以动态管理和创建。Controller Manager包含一个TokenCleaner Controller，如果到期时可以删除Bootstrap Tokens Controller。Tokens的形式是[a-z0-9]{6}.[a-z0-9]{16}。第一个组件是Token ID，第二个组件是Token Secret。可以在HTTP header中指定Token:`Authorization: Bearer 781292.db7bc3a58fc5f07e`。Bootstrap Token当前仅用作kubelet的客户端证书。

	kubelet发起的CSR请求都是由 controller manager 来做实际签署的，对于 controller manager 来说，TLS bootstrapping 下 kubelet 发起的 CSR 请求大致分为以下三种

		- nodeclient
		
			kubelet 以 O=system:nodes 和 CN=system:node:(node name) 形式发起的 CSR 请求，仅在第一次启动时会产生

		- selfnodeclient
		
			kubelet client renew自己的证书发起的CSR请求(与上一个证书就有相同的O和CN)，kubelet renew自己作为client apiserver通讯时使用的证书产生

		- selfnodeserver
		
			kubelet server renew自己的证书发起的CSR请求，kubelet首次申请或后续renew自己的10250 api端口证书时产生

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