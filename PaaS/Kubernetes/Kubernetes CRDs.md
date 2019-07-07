# Kubernetes CRDs自定义资源 #

CustomResourceDefinition（CRD）是v1.7+新增的无需改变代码就可以扩展Kubernetes API的机制，用来管理自定义对象，存在于所有namespace下。CustomResourceDefinition (CRD)是一个内建的API, 它提供了一个简单的方式来创建自定义资源。

## CRD创建流程 ##

当创建一个新的自定义资源定义（CRD）时，Kubernetes API Server通过创建一个新的RESTful资源路径进行应答。

1，定义和创建自定义资源kind: CustomResourceDefinition CRD
如下，首先需要先定义和创建一个自定义资源kind: CustomResourceDefinition，指定API Group的名称如group: networking.istio.io，

```
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    helm.sh/hook: crd-install
  creationTimestamp: 2018-10-15T10:06:40Z
  generation: 1
  labels:
    app: istio-pilot
  name: virtualservices.networking.istio.io
  resourceVersion: "226955722"
  selfLink: /apis/apiextensions.k8s.io/v1beta1/customresourcedefinitions/virtualservices.networking.istio.io
  uid: 02eb8d88-d062-11e8-a4b4-005056b84e17
spec:
  group: networking.istio.io
  names:
    kind: VirtualService
    listKind: VirtualServiceList
    plural: virtualservices
    singular: virtualservice
  scope: Namespaced
  version: v1alpha3
status:
  acceptedNames:
    kind: VirtualService
    listKind: VirtualServiceList
    plural: virtualservices
    singular: virtualservice
  conditions:
  - lastTransitionTime: 2018-10-15T10:06:40Z
    message: no conflicts found
    reason: NoConflicts
    status: "True"
    type: NamesAccepted
  - lastTransitionTime: 2018-10-15T10:06:40Z
    message: the initial names have been accepted
    reason: InitialNamesAccepted
    status: "True"
    type: Established
```
- **name**：用于定义CRD的名字，后缀需要跟group一致，前缀需要跟names中的plural一致。
- **group以及version用于标识restAPI**：即/apis//。
- **scope**: 表明作用于，可以是基于namespace的，也可以是基于集群的。 如果是基于namespace的。则API格式为：/apis/{group}/v1/namespaces/{namespace}/{spec.names.plural}/… 如果是基于cluster的。则API格式为：/apis/{group}/v1/{spec.names.plural}/… 
- **names**：描述了一些自定义资源的名字以及类型的名字（重点是plural定义以及kind定义，因为会在url或者查询资源中用的到）。

这样就创建了一个新的区分命名空间的RESTful API断点：/apis/networking.istio.io/v1alpha3/namespaces/*/virtualservices/...，然后可以使用此端点URL来创建和管理自定义对象，这些对象的kind就是上面创建的CRD中指定的kind: VirtualService对象。

2，创建一个CRD的自定义对象
在CRD对象创建完成之后就创建自定义对象(instances)了，这些自定义对象实例就可以类似Kubernetes的常用对象如Deployment、Service、Pod等一样进行CURD操作了。  自定义对象可以包含自定义的字段，这些字段可以包含任意的JSON，具体的字段要根据对象去定义，主要是spec域。

```
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: wudebao-web
spec:
  hosts:
  - "www.wudebao-web.com"
  gateways:
  - wudebao-web-gateway
  http:
  - match:
    - uri:
        exact: /wudebao
    route:
    - destination:
        port:
          number: 54321
        host: wudebao-web
```
