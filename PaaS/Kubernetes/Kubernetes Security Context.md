# Kubernetes Security Context #

Security Context的的目是限制不可信容器的行为，保护系统和其他容器不受其影响。Kubernetes提供了三种配置Security Context的方法：

- **Container-level Security Context**：仅应用到指定的容器，并且不会影响Volume。比如设置容器运行在特权模式：

	```yaml
	apiVersion: v1
	kind: Pod
	metadata:
	  name: hello-world
	spec:
	  containers:
	    - name: hello-world-container
	      # The container definition
	      # ...
	      securityContext:
	        privileged: true
	```

- **Pod-level Security Context**：应用到Pod内所有容器以及Volume(包括fsGroup和selinuxOptions)

	```yaml
	apiVersion: v1
	kind: Pod
	metadata:
	  name: hello-world
	spec:
	  containers:
	  # specification of the pod's containers
	  # ...
	  securityContext:
	    fsGroup: 1234
	    supplementalGroups: [5678]
	    seLinuxOptions:
	      level: "s0:c123,c456"
	```

- **Pod Security Policies（PSP）**：应用到集群内部所有Pod以及Volume,使用PSP需要API Server开启extensions/v1beta1/podsecuritypolicy，并且配置PodSecurityPolicyadmission控制器。

	支持的控制项如下：
	![](img/psp.png)

	示例：（限制容器的host端口范围为8000-8080）

	```yaml
	apiVersion: extensions/v1beta1
	kind: PodSecurityPolicy
	metadata:
	  name: permissive
	spec:
	  seLinux:
	    rule: RunAsAny
	  supplementalGroups:
	    rule: RunAsAny
	  runAsUser:
	    rule: RunAsAny
	  fsGroup:
	    rule: RunAsAny
	  hostPorts:
	  - min: 8000
	    max: 8080
	  volumes:
	  - '*'
	 ```


