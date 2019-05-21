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


## 准入控制 ##

摘自[kubernetes API Server 权限管理实践](https://www.cnblogs.com/fengjian2016/p/8134068.html)

准入控制admission controller本质上为一段准入代码，在对Kubernetes API的请求过程中，顺序为先经过认证、授权，然后进入准入流程，再对目标对象进行操作。这个准入代码在apiserver中，而且必须被编译到二进制文件中才能被执行。
在对集群进行请求时，每个准入控制代码都按照一定顺序执行。如果有一个准入控制拒绝了此次请求，那么整个请求的结果将会立即返回，并提示用户相应的error信息。

**Admission Controller工作流**

摘自[Kubernetes准入控制Admission Controller介绍](https://mp.weixin.qq.com/s?__biz=MzIzNzU5NTYzMA==&mid=2247484645&idx=1&sn=87b20f0bf94922cca0d0817d4be8281d&chksm=e8c77a64dfb0f37227134ec8da685f9198d814b6a223a668e55c8bb592d4ab65203742e4acfd&token=1308536270&lang=zh_CN#rd)

![](img/admission_controller.png)

API Server接收到客户端请求后首先进行认证鉴权，认证鉴权通过后才会进行后续的endpoint handler处理。

- 1)当API Server接收到对象后首先根据http的路径可以知道对象的版本号，然后将request body反序列化成versioned object.
 
- 2)versioned object转化为internal object，即没有版本的内部类型，这种资源类型是所有versioned类型的超集。只有转化为internal后才能适配所有的客户端versioned object的校验。

- 3)Admission Controller具体的admit操作，可以通过这里修改资源对象，例如为Pod挂载一个默认的Service Account等。
 
- 4)API Server internal object validation，校验某个资源对象数据和格式是否合法，例如：Service Name的字符个数不能超过63等。
 
- 5)Admission Controller validate，可以自定义任何的对象校验规则。
 
- 6)internal object转化为versioned object，并且持久化存储到etcd。

Kubernetes 1.10之前版本在APIerver中配置`--admission_control`参数可以进行准入控制的配置，它的值为一串用逗号连接的、有序的准入模块列表。Kubernetes 1.10之后的版本，`--admission-control`已经废弃，建议使用`--enable-admission-plugins` `--disable-admission-plugins`指定需要打开或者关闭的Admission Controller。 同时用户指定的顺序并不影响实际Admission Controllers的执行顺序，对用户来讲非常友好。

它的模块如下：

- AlwaysAdmit：允许所有请求
 
- AlwaysDeny：禁止所有请求，多用于测试环境。

- AlwaysPullImages: 该插件修改每个新的Pod，强制pull最新镜像，这在多租户群集中非常有用，以便私有镜像只能由拥有授权凭据的用户使用。
 
- DenyExecOnPrivileged(已弃用)：它会拦截所有想在privileged container上执行命令的请求。如果自己的集群支持privileged container，自己又希望限制用户在这些privileged container上执行命令，那么强烈推荐使用它。此功能已合并到`DenyEscalatingExec`中。

- DenyEscalatingExec: 禁止privileged container的exec和attach操作。

- ImagePolicyWebhook: 通过webhook决定image策略，需要同时配置`--admission-control=ImagePolicyWebhook`。

	**配置文件格式**
	
	ImagePolicyWebhook使用admission配置文件`--admission-control-config-file`为Backend Behavior设置配置选项。该文件可以是json或yaml并具有以下格式:

	```
	{
	  "imagePolicy": {
	     "kubeConfigFile": "path/to/kubeconfig/for/backend",
	     "allowTTL": 50,           // time in s to cache approval
	     "denyTTL": 50,            // time in s to cache denial
	     "retryBackoff": 500,      // time in ms to wait between retries
	     "defaultAllow": true      // determines behavior if the webhook backend fails
	  }
	}
	```

 
- ServiceAccount：这个plug-in将serviceAccounts实现了自动化，如果想要使用ServiceAccount对象，那么强烈推荐使用它。

	一个serviceAccount为运行在pod内的进程添加了相应的认证信息。当准入模块中开启了此插件（默认开启），那么当pod创建或修改时他会做一下事情：

	  -	如果pod没有serviceAccount属性，将这个pod的serviceAccount属性设为“default”；

      - 确保pod使用的serviceAccount始终存在；

      - 如果LimitSecretReferences 设置为true，当这个pod引用了Secret对象却没引用ServiceAccount对象，弃置这个pod；

      - 如果这个pod没有包含任何ImagePullSecrets，则serviceAccount的ImagePullSecrets被添加给这个pod；

      - 如果MountServiceAccountToken为true，则将pod中的container添加一个VolumeMount 。

- SecurityContextDeny：这个插件将会将使用了SecurityContext的pod中定义的选项全部失效。SecurityContext在container中定义了操作系统级别的安全设定（uid, gid, capabilities, SELinux等等）。

- ResourceQuota：它会观察所有的请求，确保在namespace中ResourceQuota对象处列举的container没有任何异常。 如果在kubernetes中使用了ResourceQuota对象，就必须使用这个插件来约束container。推荐在admission control参数列表中，这个插件排最后一个。

- LimitRanger：它会观察所有的请求，确保没有违反LimitRanger对象中枚举的任何限制Namespace，如果在Kubernetes Deployment中使用了LimitRanger对象，则必须使用此插件

- InitialResources： 根据镜像的历史使用记录，为容器设置默认资源请求和limits。

- NamespaceLifecycle： 该插件确保处于Termination状态的Namespace不再接收新的对象创建请求，并拒绝请求不存在的Namespace。该插件还可以防止删除系统保留的Namespace:`default，kube-system，kube-public`。

- DefaultStorageClass: 该插件将观察PersistentVolumeClaim，并自动设置默认的Storage Class。当没有配置默认Storage Class时，此插件不会执行任何操作。当有多个Storage Class被标记为默认值时，它也将拒绝任何创建，管理员必须重新访问StorageClass对象，并且只标记一个作为默认值。此插件不用于PersistentVolumeClaim的更新，仅用于创建。

- DefaultTolerationSeconds: 该插件设置Pod的默认forgiveness toleration为5分钟。

- PodSecurityPolicy: 该插件用于创建和修改pod，使用Pod Security Policies时需要开启。

- NodeRestriction: 此插件限制kubelet修改Node和Pod对象，这样的kubelets只允许修改绑定到Node的Pod API对象，以后版本可能会增加额外的限制。

- NamespaceExists：它会观察所有的请求，如果请求尝试创建一个不存在的namespace，则这个请求被拒绝。



**推荐的插件配置：**

- Kubernetes > = 1.6.0

	`--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds`

- Kubernetes > = 1.4.0

	`--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota`

- Kubernetes > = 1.2.0

	`--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,ResourceQuota`

- Kubernetes > = 1.0.0

	`--admission-control=NamespaceLifecycle,LimitRanger,SecurityContextDeny,ServiceAccount,PersistentVolumeLabel,ResourceQuota`


## RBAC ##

RBAC（Role-Based Access Control），允许通过Kubernetes API动态配置授权策略。

RBAC被映射为四种Kubernetes顶级资源对象：

- 角色（Role）：Role和ClusterRole
	- 角色表示一组权限的规则，累积规则
	- Role适用带namespace的资源，ClusterRole适用集群资源或非资源API

- 建立用户与角色的绑定关系：RoleBinding和ClusterRoleBinding
	- RoleBinding和ClusterRoleBinding的区别在于是否是namespace的资源
	- 角色绑定包含了一组相关主体（即subject，包括用户、用户组或者serviceAccount）以及对呗授予角色的引用	