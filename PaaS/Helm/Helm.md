# Helm #

[https://hub.kubeapps.com/](https://hub.kubeapps.com/ "公共仓库")

## 架构 ##

![](img/helm01.png)

**Helm客户端**：一个供终端用户使用的命令行工具，客户端负责如下的工作：

- 本地chart开发
- 管理仓库
- 与Tiller服务器交互
- 发送需要被安装的charts
- 请求关于发布版本的信息
- 请求更新或者卸载已安装的发布版本

**Tiller服务器**： Tiller服务部署在Kubernetes集群中，Helm客户端通过与Tiller服务器进行交互，并最终与Kubernetes API服务器进行交互。 Tiller服务器负责如下的工作：

- 监听来自于Helm客户端的请求
- 组合chart和配置来构建一个发布
- 在Kubernetes中安装，并跟踪后续的发布
- 通过与Kubernetes交互，更新或者chart
- 客户端负责管理chart，服务器发展管理发布。

## Helm Hook ##

hook的任务执行的时候，Tiller会阻塞，任务执行失败，则应用发布失败。Hook创建的资源不作为release的一部分进行跟踪或管理。一旦Tiller验证hook已经达到其就绪状态，就不再跟踪它了，即无法再对其进行操作。所以hook创建的资源，无法使用`helm delete`删除。

**支持Hook类型**

- pre-install：在模板文件被渲染之后、而在Kubernetes创建任何资源创建之前执行。
- post-install：在Kubernetes加载全部的资源之后执行。
- pre-delete：在Kubernetes删除任何resource之前执行。
- post-delete：在一个release的全部资源被删除之后执行。
- pre-upgrade：在模板渲染之后，而在Kubernetes加载任何资源之前执行。
- post-upgrade：在Kubernetes更新完全部resource之后执行。
- pre-rollback：在模板被渲染之后、而在Kubernetes执行对全部resource的回滚之前。
- post-rollback：在Kubernetes的全部resource已经被修改之后执行。

**Hook的权重**

权重的执行顺序：由负到正、从小到大
默认权重：默认为0

**hook的删除策略**

- “hook-succeeded” ：指定Tiller应该在hook成功执行后删除hook。
- “hook-failed” ：指定如果hook在执行期间失败，Tiller应该删除hook。
- “before-hook-creation” ： 指定Tiller应在删除新hook之前删除以前的hook。

**Hook设置**

```
  annotations:
    # This is what defines this resource as a hook. Without this line, the
    # job is considered part of the release.
    "helm.sh/hook": post-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
```

## Chart ##

**Chart结构**

```
chart                                # Chart的名字，也就是目录的名字
├── charts                           # Chart所依赖的子Chart
│   ├── jenkins-0.14.0.tgz           # 被“mychart”依赖的其中一个subChart
│   ├── mysubchart                   # 被“mychart”依赖的另一个subChart
│   │   ├── charts                   
│   │   ├── Chart.yaml
│   │   ├── templates
│   │   │   └── configmap.yaml
│   │   └── values.yaml
│   └── redis-1.1.17.tgz             
├── Chart.yaml                       # 记录关于该Chart的描述信息：比如名称、版本等等
├── config1.toml                     # 其他文件：可以是任何文件
├── config2.toml                     # 其他文件：可以是任何文件
├── requirements.yaml                # 记录该Chart的依赖，类似pom.xml
├── templates                        # 存放模版文件，模板也就是将k8s的yml文件参数化，最终还是会被helm处理成k8s的正常yml文件，然后用来部署对应的资源
│   ├── configmap.yaml               # 一个ConfigMap资源模版
│   ├── _helpers.tpl                 # "_"开头的文件不会本部署到k8s上，可以用于定于通用信息，在其他地方应用，如loables
│   └── NOTES.txt                    # 在执行helm instll安装此Chart之后会被输出到屏幕的一些自定义信息
└── values.yaml                      # 参数定义文件，这里定义的参数最终会应用到模版中
```

**内置变量**

- **Release.Name**: The name of the release (not the chart)
- **Release.Tim**e: The time the chart release was last updated. This will match the Last Released time on a Release object.
- **Release.Namespace**: The namespace the chart was released to.
- **Release.Service**: The service that conducted the release. Usually this is Tiller.
- **Release.IsUpgrade**: This is set to true if the current operation is an upgrade or rollback.
- **Release.IsInstall**: This is set to true if the current operation is an install.
- **Release.Revision**: The revision number. It begins at 1, and increments with each helm upgrade.
- **Chart**: The contents of the Chart.yaml. Thus, the chart version is obtainable as Chart.Version and the maintainers are in Chart.Maintainers.
- **Files**: A map-like object containing all non-special files in the chart. This will not give you access to templates, but will give you access to additional files that are present (unless they are excluded using .helmignore). Files can be accessed using {{index .Files "file.name"}} or using the {{.Files.Get name}} or {{.Files.GetString name}} functions. You can also access the contents of the file as []byte using {{.Files.GetBytes}}
- **Capabilities**: A map-like object that contains information about the versions of Kubernetes ({{.Capabilities.KubeVersion}}, Tiller ({{.Capabilities.TillerVersion}}, and the supported Kubernetes API versions ({{.Capabilities.APIVersions.Has "batch/v1")

注意：其中`Files`变量会经常用到，configmap中一些内容可能与helm的语法有冲突，在部署时并不想helm去渲染这部分内容，可以将这部分内容放到`config.toml`这种文件中，然后在configmap中使用`{{.Files.Get config.toml}}`来获取。
