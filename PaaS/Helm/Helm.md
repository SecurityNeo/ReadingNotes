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
- crd-install：在运行其他检查之前添加CRD资源，只能用于chart中其他的资源清单定义的 CRD 资源。

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

- **Release.Name**: release名称
- **Release.Time**: release的最近更新时间
- **Release.Namespace**: release的namespace
- **Release.Service**: release服务的名称（始终是Tiller）
- **Release.IsUpgrade**: 如果当前操作是升级或回滚，则将其设置为true
- **Release.IsInstall**: 如果当前操作是安装，则设置为true
- **Release.Revision**: 此release的修订版本号。它从1开始，每helm upgrade一次增加一个
- **Chart**: Chart.yaml文件的内容。任何数据Chart.yaml将在这里访问。
- **Files**:提供对chart中所有非特殊文件的访问。虽然无法使用它来访问模板，但可以使用它来访问chart中的其他文件。
	- `Files.Get`是一个按名称获取文件的函数，例如（.Files.Get config.ini）
	- `Files.GetBytes`是将文件内容作为字节数组而不是字符串获取的函数。这对于像图片这样的东西很有用。
- **Capabilities**: 这提供了关于Kubernetes集群支持的功能的信息。
	- `Capabilities.APIVersions`是一组版本信息。
	- `Capabilities.APIVersions.Has $version`指示是否在群集上启用版本（batch/v1）。
	- `Capabilities.KubeVersion`提供了查找Kubernetes版本的方法。它具有以下值：Major，Minor，GitVersion，GitCommit，GitTreeState，BuildDate，GoVersion，Compiler和Platform。
	- `Capabilities.TillerVersion`提供了查找Tiller版本的方法。它具有以下值：SemVer，GitCommit，和 GitTreeState。
- **Template**：包含有关正在执行的当前模板的信息
- **Name**：到当前模板的namespace文件路径（例如mychart/templates/mytemplate.yaml）
- **BasePath**：当前chart模板目录的namespace路径（例如 mychart/templates）。

注意：其中`Files`变量会经常用到，configmap中一些内容可能与helm的语法有冲突，在部署时并不想helm去渲染这部分内容，可以将这部分内容放到`config.toml`这种文件中，然后在configmap中使用`{{.Files.Get config.toml}}`来获取。


## 函数 ##

[https://godoc.org/text/template](https://godoc.org/text/template)

**default函数**

允许在模板内部指定默认值，以防该值被省略。
```
drink: {{.Values.favorite.drink | default "tea" | quote}}
```

**运算符函数**

在模版中有一些运算符（eq，ne，lt，gt，and，or等等）

```
{{if and .Values.fooString (eq .Values.fooString "foo") }}
    {{...}}
{{end}}
```

## 控制结构 ##

Helm的模板语言提供了以下控制结构：

- if/else 用于创建条件块
- with 指定范围
- range，它提供了一个 “for each” 风格的循环

除此之外，它还提供了一些声明和使用命名模板段的操作：

- define 在模板中声明一个新的命名模板
- template 导入一个命名模板
- block 声明了一种特殊的可填写模板区域

**if/else**

```
{{if PIPELINE}}
  # Do something
{{else if OTHER PIPELINE}}
  # Do something else
{{else}}
  # Default case
{{end}}
```

如果值为如下情况，则管道评估为false:

- 一个布尔型的假
- 一个数字零
- 一个空的字符串
- 一个 nil（空或 null）
- 一个空的集合（map，slice，tuple，dict，array）

**控制空格**

- 使用特殊字符修改模板声明的大括号语法，以告诉模板引擎填充空白。{{-（添加了破折号和空格）表示应该将空白左移，而 -}} 意味着应该删除右空格。注意！换行符也是空格！
- 使用indent函数（{{indent 2 "mug:true"}}）

**使用`with`修改范围**

`with`控制着变量作用域。语法：
```
{{with PIPELINE}}
  # restricted scope
{{end}}
```
示例：
```
 myvalue: "Hello World"
  {{- with .Values.favorite}}
  drink: {{.drink | default "tea" | quote}}
  food: {{.food | upper | quote}}
  {{- end}}
```
引用`.drink`和`.food`无需对其进行限定。这是因为该`with`声明设置 . 为指向`.Values.favorite`。在`{{end}}`后 . 复位其先前的范围。

**循环`range`**

示例：

```
favorite:
  drink: coffee
  food: pizza
pizzaToppings:
  - mushrooms
  - cheese
  - peppers
  - onions
```
----------

```
  myvalue: "Hello World"
  {{- with .Values.favorite}}
  drink: {{.drink | default "tea" | quote}}
  food: {{.food | upper | quote}}
  {{- end}}
  toppings: |-
    {{- range .Values.pizzaToppings}}
    - {{. | title | quote}}
    {{- end}}
```

----------
渲染结果：
```
  myvalue: "Hello World"
  drink: "coffee"
  food: "PIZZA"
  toppings: |-
    - "Mushrooms"
    - "Cheese"
    - "Peppers"
    - "Onions"
```
注：YAML中的`|-`标记表示一个多行字符串。

## 命名模板 ##

**用`define`和`template`声明和使用模板**

语法：
```
{{ define "MY.NAME" }}
  # body of template here
{{ end }}
```
示例：
```
{{- define "mychart.labels" }}
  labels:
    generator: helm
    date: {{ now | htmlDate }}
{{- end }}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
  {{- template "mychart.labels" }}
data:
  myvalue: "Hello World"
  {{- range $key, $val := .Values.favorite }}
  {{ $key }}: {{ $val | quote }}
  {{- end }}
```

Helm chart 通常将这些模板放入partials文件中，通常是_helpers.tpl。按照惯例，define函数应该有一个简单的文档块`（{{/* ... */}}）`来描述他们所做的事情。


## 文件访问 ##

**Glob模式**

`.Glob`返回一个Files类型，可以调用Files返回对象的任何方法。

示例：
文件结构
```
foo/:
  foo.txt foo.yaml
bar/:
  bar.go bar.conf baz.yaml
```

```
{{$root := .}}
{{range $path, $bytes := .Files.Glob "**.yaml"}}
{{$path}}: |-
{{$root.Files.Get $path}}
{{end}}
```

```
{{range $path, $bytes := .Files.Glob "foo/*"}}
{{$path.base}}: '{{ $root.Files.Get $path | b64enc }}'
{{end}}
```

**ConfigMap和Secrets工具函数**

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: conf
data:
  {{- (.Files.Glob "foo/*").AsConfig | nindent 2 }}
---
apiVersion: v1
kind: Secret
metadata:
  name: very-secret
type: Opaque
data:
  {{(.Files.Glob "bar/*").AsSecrets | nindent 2 }}
```

**访问文件的每一行**

```
data:
  some-file.txt: {{range .Files.Lines "foo/bar.txt"}}
    {{.}}{{ end }}
```

## 自定义资源定义(CRD) ##

对于CRD，声明必须在该CRDs种类的任何资源可以使用之前进行注册。注册过程有时需要几秒钟。
两种方法：

- 独立的chart

	将CRD定义放在一个chart中，然后将所有使用该CRD的资源放入另一个chart中。在这种方法中，每个chart必须单独安装。

- 预安装hook

	在CRD定义中添加一个`crd-install`钩子，以便在执行chart的其余部分之前完全安装它。

注意：如果使用crd-install hook创建CRD ，则该CRD定义在helm delete运行时不会被删除。

## 一些常用技巧 ##

**删除Helm release时保留一些资源**

```
metadata:
  annotations:
    "helm.sh/resource-policy": keep
```
注释`"helm.sh/resource-policy": keep`指示Tiller在`helm delete`操作过程中跳过此资源。但是，此资源变成孤儿资源。Helm将不再以任何方式管理它。