# KrakenD #

[https://www.krakend.io/](https://www.krakend.io/)

[https://github.com/devopsfaith](https://github.com/devopsfaith)

- 配置解析
	- Service
		- [TLS Endpoint](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#tls-endpoint)

	- Endpoint
		- [Endpoint Rate Limiting](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#endpoint-rate-limiting)
		- [Response Manipulation](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#response-manipulation)
		- [参数与头部转发](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#%E5%8F%82%E6%95%B0%E4%B8%8E%E5%A4%B4%E9%83%A8%E8%BD%AC%E5%8F%91)
		- [Content Types](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#content-types)
		- [Backend的顺序调用](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#backend%E7%9A%84%E9%A1%BA%E5%BA%8F%E8%B0%83%E7%94%A8)
		- [Static Proxy](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#static-proxy)

	- Backend
		- [断路器](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#%E6%96%AD%E8%B7%AF%E5%99%A8)
		- [并发请求](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#%E5%B9%B6%E5%8F%91%E8%AF%B7%E6%B1%82)
		- [数据缓存](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#%E6%95%B0%E6%8D%AE%E7%BC%93%E5%AD%98)
		- [Traffic Shadowing](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#traffic-shadowing)
		- [flatmap](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#flatmap)
	- 认证
		- [JWT Validation](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#jwt-validation)
		- [JWT Signing](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#jwt-signin)
		- [OAuth 2.0支持](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#oauth-20%E6%94%AF%E6%8C%81)
	- 服务发现
		- [静态](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#%E9%9D%99%E6%80%81)
		- [DNS SRV](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#dns-srv)
		- [ETCD](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#etcd)
		- [Eureka](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#eureka)
	- 限速与节流
		- [超时](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#%E8%B6%85%E6%97%B6)
		- [控制闲置连接](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#%E6%8E%A7%E5%88%B6%E9%97%B2%E7%BD%AE%E8%BF%9E%E6%8E%A5)
	- Logging和Metrics
		- [Logging](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#logging)
		- [Metrics](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Microservice/KrakenD.md#metrics)

KrakenD是一个收费的API网关生成器和代理生成器，位于客户端和所有源服务器之间，添加了一个新层，可以消除客户机的所有复杂性，为客户机只提供UI需要的信息。 KrakenD充当了许多源的集合，可以将许多源集成到单个端点中，并允许对响应进行分组、包装。 另外它支持大量middelwares和插件，允许扩展功能，比如添加Oauth授权或者安全层。

## 配置解析 ##

### Service ###

#### TLS Endpoint ####

KrakenD需要在全局配置中开启TLS，一旦开启了TLS，KrakenD将不会响应任何HTTP请求。

示例：

```
{
  "version": 2,
  "tls": {
    "public_key": "/path/to/cert.pem",
    "private_key": "/path/to/key.pem"
  }
}
```

- public_key: 公钥文件，绝对或相对路径（相对工作目录）
- private_key：私钥文件，绝对或相对路径（相对工作目录）

可选配置：

- disabled (boolean): 临时关闭TLS，用户开发测试环境
- min_version (string): 最低TLS版本 (SSL3.0, TLS10, TLS11 或者 TLS12)
- max_version (string): 最高TLS版本 (SSL3.0, TLS10, TLS11 或者 TLS12)
- curve_preferences (integer array): 椭圆数字签名算法签名长度(23表示CurveP256, 24表示CurveP384，25表示CurveP521)
- prefer_server_cipher_suites (boolean): 强制使用服务器提供的密码套件，而不是使用客户端建议的套件。
- cipher_suites (integer array): 加密算法，可选算法如下：
	- 5: TLS_RSA_WITH_RC4_128_SHA
	- 10: TLS_RSA_WITH_3DES_EDE_CBC_SHA
	- 47: TLS_RSA_WITH_AES_128_CBC_SHA
	- 53: TLS_RSA_WITH_AES_256_CBC_SHA
	- 60: TLS_RSA_WITH_AES_128_CBC_SHA256
	- 156: TLS_RSA_WITH_AES_128_GCM_SHA256
	- 157: TLS_RSA_WITH_AES_256_GCM_SHA384
	- 49159: TLS_ECDHE_ECDSA_WITH_RC4_128_SHA
	- 49161: TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
	- 49162: TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
	- 49169: TLS_ECDHE_RSA_WITH_RC4_128_SHA
	- 49170: TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
	- 49171: TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
	- 49172: TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
	- 49187: TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
	- 49191: TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
	- 49199: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	- 49195: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
	- 49200: TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
	- 49196: TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
	- 52392: TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
	- 52393: TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305


### Endpoint ###

#### Endpoint Rate Limiting ####

可以限制Endpoint的速率，有两种限制方法，一种为限制每个Endpoint的速率，一种为限制每个Endpoint可接受单个客户端的速率。
注意： 开始针对每个客户端的速率限制将在很大程度上影响API网关的性能！

示例：

```
{
    "version": 2,
    "endpoints": [
      {
          "endpoint": "/happy-hour",
          "extra_config": {
              "github.com/devopsfaith/krakend-ratelimit/juju/router": {
                  "maxRate": 0,
                  "clientMaxRate": 0
              }
          }
          ...
      },
      {
          "endpoint": "/limited-endpoint",
          "extra_config": {
            "github.com/devopsfaith/krakend-ratelimit/juju/router": {
                "maxRate": 50,
                "clientMaxRate": 5,
                "strategy": "ip"
              }
          },
          ...
      },
      {
          "endpoint": "/user-limited-endpoint",
          "extra_config": {
            "github.com/devopsfaith/krakend-ratelimit/juju/router": {
                "clientMaxRate": 10,
                "strategy": "header",
                "key": "X-Auth-Token"
              }
          },
          ...
      }
```

- maxRate: 不配置或配置`"maxRate": 0 `表示不进行限速。当达到Endpoint的限速配置后，客户端将收到`503 Service Unavailable`。
- clientMaxRate：不配置或配置`"clientMaxRate": 0`表示不进行限速。当达到Endpoint的限速配置后，客户端将收到`429 Too Many Requests`。
	- "strategy": "ip"：通过客户端IP地址识别是否为同一个客户端。
	- "strategy": "header"： 通过请求包对应的头部来识别是否为同一个客户端。


#### Response Manipulation ####

KrakenD可以修改Backend返回的数据，主要包含以下几类：Merging（数据聚合）、Filtering（数据过滤）、Grouping（数据分组）、Mapping (key修改)、Target（数据截取）、Collection（数据组合）。

**Merging**

将所有Backend返回的数据聚合到一个字典中，如果多个Backend中有key冲突，则对应key的值将取返回数据最慢的那个Backend的值。

示例：
```
"endpoints": [
    {
      "endpoint": "/abc",
      "timeout": "800ms",
      "method": "GET",
      "backend": [
        {
          "url_pattern": "/a",
          "encoding": "json",
          "host": [
            "http://service-a.company.com"
          ]
        },
        {
          "url_pattern": "/b",
          "encoding": "xml",
          "host": [
            "http://service-b.company.com"
          ]
        }
      ]
    }
```

- timeout: KrakenD不会一直等待Backend返回数据，需要为其配置一个超时时间，此配置可以为全局配置，也可以为某个具体的Backend配置超时时间，如果都有配置，则以对应Backend中的超时时间为准。如果达到超时时间，部分Backend仍然没有返回数据，KrakenD会返回达到超时时间前获取到的部分数据，在头部字段中不会返回cache头部，并且新增一个头部：`x-krakend-completed: false`，反之该头部为`x-krakend-completed: true`。

**Filtering**

KrakenD可以过滤Backend返回的数据，这样可以极大地降低客户端带宽的占用，实现方式有两种：`Blacklist`和`Whitelist`。
注意：`Blacklist`和`Whitelist`只能取其中之一，因为它们本身就是两个冲突的行为，另外从性能上来讲，`Blacklist`要比`Whitelist`快。

- Blacklist： 指定哪些字段内容不转发到客户端，可以支持嵌套，但不支持数组，关于数组可以参考`flatmap_filter`。

	示例：
	```
	{
	  "endpoint": "/posts/{user}",
	  "method": "GET",
	  "backend": [
	    {
	      "url_pattern": "/posts/{user}",
	      "host": [
	        "https://jsonplaceholder.typicode.com"
	      ],
	      "blacklist": [
	        "body",
	        "user.userId"
	      ]
	    }
	  ]
	}
	```

- Whitelist： 指定只转发这部分字段内容到客户端，可以支持嵌套，但不支持数组，关于数组可以参考`flatmap_filter`。

	示例：
	
	```
	{
	  "endpoint": "/posts/{user}",
	  "method": "GET",
	  "backend": [
	    {
	      "url_pattern": "/posts/{user}",
	      "host": [
	        "https://jsonplaceholder.typicode.com"
	      ],
	      "whitelist": [
	        "id",
	        "title"
	      ]
	    }
	  ]
	}
	```

**Grouping**

KrakenD可以为Backend返回的数据分组，简单来说，可以为某个Backend指定gourp，转发给客户端时，此Backend的Response将嵌套到指定的group中，这在一定程度上解决多个Backend返回了相同的key，但key的含义不一样，都是客户端需要的这种情况。
注意：同一个Endpoint的每个Backend中group应该是不同的，如果相同，那最终转发到客户端的数据中对应group的值为最慢的那个Backend返回的数据。

示例：

```
{
  "endpoint": "/users/{user}",
  "method": "GET",
  "backend": [
    {
      "url_pattern": "/users/{user}",
      "host": [
        "https://jsonplaceholder.typicode.com"
      ]
    },
    {
      "url_pattern": "/posts/{user}",
      "host": [
        "https://jsonplaceholder.typicode.com"
      ],
      "group": "last_post"
    }
  ]
}
```
返回数据：
```
{
  "id": 1,
  "phone": "1-770-736-8031 x56442",
  "website": "hildegard.org"
  ...
  "last_post": {
    "id": 1,
    "userId": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit"
  }
}
```

**Mapping**

KrakenD可以对Backend返回的字段进行重命名。目前貌似只修改支持最外层的key。

示例：

```
{
  "endpoint": "/users/{user}",
  "method": "GET",
  "backend": [
    {
      "url_pattern": "/users/{user}",
      "host": [
        "https://jsonplaceholder.typicode.com"
      ],
      "mapping": {
        "email": "personal_email"
      }
    }
  ]
}
```
返回数据：
```
{
  "id": 1,
  "name": "Leanne Graham",
  "username": "Bret",
  "personal_email": "Sincere@april.biz",
  ...
}
```

**Target**

KrakenD能获取Backend某个字段的内容，并只将这个字段的内容转发给客户端。
注意，`Target`跟`Whitelist`不一样。`Target`可以与`whitelist`、`mapping`等混用，`Target`发生在这些动作之后。

示例：
	原始数据:
	```
	{
	  "apiVersion":"2.0",
	  "data": {
	    "updated":"2010-01-07T19:58:42.949Z",
	    "totalItems":800,
	    "startIndex":1,
	    "itemsPerPage":1,
	    "items":[]
	  }
	}
	```

	配置：
	```
	{
	  "endpoint": "/foo",
	  "method": "GET",
	  "backend": [
	    {
	      "url_pattern": "/bar",
	      "target": "data"
	    }
	  ]
	}
	```

	转发的数据内容：
	```
	{
	    "updated":"2010-01-07T19:58:42.949Z",
	    "totalItems":800,
	    "startIndex":1,
	    "itemsPerPage":1,
	    "items":[]
	}
	```

**Collection**

KrakenD期望Backend返回的数据都是一个字典，当Backend返回的数据是一个数组时，添加配置`"is_collection": true`，KrakenD会将返回的字段转换为一个字典，默认的key为`collection`，可以通过`mapping`自定义key值。
示例：
```
"endpoints": [
    {
      "endpoint": "/posts",
      "backend": [
        {
          "url_pattern": "/posts",
          "host": ["http://jsonplaceholder.typicode.com"],
          "sd": "static",
          "is_collection": true,
          "mapping": {
            "collection": "myposts"
          }
        }
      ]
    }
]
```


#### 参数与头部转发 ####

默认情况下，为了安全考虑，KrakenD不转发客户端发送的任何参数和头部字段，如果有类似需求，需要做相关配置。

- 参数转发

	在Endpoint的`querystring_params`中可指定允许转发的参数key，相当于一个白名单。
	示例，KrakenD只会转发参数a和b，其它的都会被忽略：
	```
	{
	  "version": 2,
	  "endpoints": [
	    {
	      "endpoint": "/v1/foo",
	      "querystring_params": [
	        "a",
	        "b"
	      ],
	      "backend": [
	        {
	          "url_pattern": "/catalog",
	          "host": [
	            "http://some.api.com:9000"
	          ]
	        }
	      ]
	    }
	  ]
	}
	```
	如果需要配置所有参数转发，可使用通配符`*`，如下(不建议这样做，有可能会造成注入漏洞风险)：
	```
	"querystring_params":[  
	      "*"
	]
	```

	另外，可以通过配置变量的方式来强制客户端传递一些参数，如下：
	```
	{
	        "endpoint": "/v3/{channel}/foo",
	        "backend": [
	                {
	                        "host": ["http://backend"],
	                        "url_pattern": "/foo?channel={channel}"
	                }
	        ]
	}
	```
	注意：变量配置和`querystring_params`结合使用的情况，看下边的例子。

	配置文件：
	```
	{
	        "endpoint": "/v3/{channel}/foo",
	        "querystring_params": [
	                "page",
	                "limit"
	        ],
	        "backend": [
	                {
	                        "host": ["http://backend"],
	                        "url_pattern": "/foo?channel={channel}"
	                }
	        ]
	}
	```
	客户端请求： `http://krakend/v3/iOS/foo?limit=10&evil=here`，Backend收到的请求是`/foo?limit=10`;
	客户端请求： `http://krakend/v3/iOS/foo?evil=here`，Backend收到的请求是`/foo?channel=foo`。

- 头部转发

	KrakenD会向Backend发送一些基本的头部字段，如下：
	```
	Accept-Encoding: gzip
	Host: localhost:8080
	User-Agent: KrakenD Version 0.9.0
	X-Forwarded-For: ::1
	```
	同样的，可以通过`headers_to_pass`来配置KrakenD转发哪些客户端的头部。也可以使用通配符`*`来配置其转发所有头部字段。
	```
	"headers_to_pass":[  
	      "*"
	]
	```

#### Content Types ####

KrakenD可以通过配置`output_encoding`来指定返回给客户端的数据格式，支持的格式有：
- json: 返回Json格式数据
- negotiate： 客户端通过`Accept`头部来决定将数据解析为何种格式，支持的格式为JSON、XML、RSS、YAML。
- string： 返回字符串格式数据
- no-op：不进行编码解码，相当于KrakenD仅仅当做一个proxy。注意：设置KrakenD为no-op模式下时，KrakenD只会将请求转发到其中一个Backend，不会做merge、Filtering、Grouping、Mapping等操作。

官方对no-op模式的解释：
 
- The KrakenD endpoint works just like a regular proxy
- The router pipe functionalities are available (e.g., rate limiting the endpoint)
- The proxy pipe functionalities are disabled (aggregate/merge, filter, manipulations, body inspection, concurrency…)
- Headers passing to the backend still need to be declared under headers_to_pass, as they hit the router layer first.
- Backend response and headers remain unchanged (including status codes)
- The body cannot be changed and is set solely by the backend
- 1:1 relationship between endpoint-backend (one backend per endpoint).

设置KrakenD为no-op模式：

- 在endpoint中添加`"output_encoding": "no-op"`
- 在Backend中添加`"encoding": "no-op"`

#### Backend的顺序调用 ####

有时候会有这样一个场景，后面一个API请求的请求体需要使用前面一个API请求的部分返回值，可以通过配置` "sequential": true`来打开此功能。
示例：
```
"endpoint": "/hotels/{id}",
"backend": [
    { <--- Index 0
        "host": [
            "https://hotels.api"
        ],
        "url_pattern": "/hotels/{id}"
    },
    { <--- Index 1
        "host": [
            "https://hotels.api"
        ],
        "url_pattern": "/destinations/{resp0_destination_id}"
    }
],
"extra_config": {
    "github.com/devopsfaith/krakend/proxy": {
        "sequential": true
    }
}
```

#### Static Proxy ####

KrakenD可以根据Backend返回数据成功与否、数据是否完成聚合来为转发到客户端的数据增加一些静态内容。此过程是在所有Backend数据聚合之后，所以请注意添加的数据不要与实际数据产生冲突，否则会被重写。支持的策略如下：
- always: 不管何种情况都在返回的数据中添加静态内容。
- success: 当所有Backend都正常返回数据的情况下在返回的数据中添加静态内容。
- complete: 当所有Backend都正常返回数据，并且数据正常聚合的情况下在返回的数据中添加静态内容。
- errored: 当Backend出现异常并返回明确的错误时在返回的数据中添加静态内容。
- incomplete: 当部分Backend没有返回数据（超时或者其它原因引起的）时在返回的数据中添加静态内容。

示例：
```
"extra_config": {
    "github.com/devopsfaith/krakend/proxy": {
        "static": {
            "strategy": "errored",
            "data": {
                // YOUR STATIC JSON OBJECT GOES HERE
            }
        }
    }
}
```


### Backend ###

#### 断路器 ####

KrakenD中有一个组件可以实现当Backend超限返回错误时中断到其上的连接，简单地可以理解为一个断路器。

**工作原理：**
首先在Backend上配置这样三个参数`interval`、`timeout`、`maxErrors`，在`interval`间隔时间内连续收到Backend返回`maxErrors`个错误，KrakenD将会中断到此Backend的所有连接，中断时间间隔为`timeout`，到达中断时长之后，KrakenD会向此Backend发送一个请求，如果这个请求仍然收到错误，KrakenD将继续中断连接`timeout`时长，如果请求正常返回，KrakenD将恢复到正常状态。

断路器的状态机：
![](img/KrakenD_Circuit_Breaker_States.png)

- CLOSED: 初始化状态，KrakenD正常转发请求到Backend。
- OPEN: KrakenD收到Backend返回所配置的错误数量后进入到OPEN状态，并且不会再转发任何请求到Backend上，此状态持续`timeout`秒。
- HALF-OPEN:`timeout`秒之后，KrakenD再发送一个请求到Backend，如果正常将进入CLOSE状态，正常转发请求，如果继续收到ERROR则重新进入OPEN状态。

**配置**

```
"endpoints": [
{
    "endpoint": "/myendpoint",
    "method": "GET",
    "backend": [
    {
        "host": [
            "http://127.0.0.1:8080"
        ],
        "url_pattern": "/mybackend-endpoint",
        "extra_config": {
            "github.com/devopsfaith/krakend-circuitbreaker/gobreaker": {
                "interval": 60,
                "timeout": 10,
                "maxErrors": 1,
                "logStatusChange": true
            }
        }
    }
    ]
```

#### 并发请求 ####

如果有多个Backend提供相同的服务，但是各个Backend的性能并不相同，这时可以为一个Endpoint的Backend配置多个主机，并且指定并发访问量，KrakenD会同时向这些主机发送所配置并发数量的API请求，并且在收到第一个主机返回的数据时就转发给客户端，然后KrakenD会忽略剩下的返回数据。所有的请求都Error之后，客户端将收到Error返回。
官方有一个示例来验证并发请求特性对整个API请求性能的影响，[点击我查看](https://www.krakend.io/docs/backends/concurrent-requests/)。

注意：需要确保这种API请求都是幂等的。

示例：
```
"endpoints": [
{
  "endpoint": "/products",
  "method": "GET",
  "concurrent_calls": 3,
  "backend": [
    {
        "host": [
            "http://server-01.api.com:8000",
            "http://server-02.api.com:8000"
        ],
        "url_pattern": "/foo",
```
注意：`concurrent_calls`与`host`的数量没有必然的关系。


#### 数据缓存 ####

作为一个API网关，为了提升响应速度，KrakenD也有缓存数据的功能，配置也相对简单，只需在Backend中开启`httpcache`中间件即可，开启之后，所有与Backend连接返回的数据都将缓存在KrakenD内存中，缓存过期时间为从Backend响应头`Cache-Control`中指定的时间。响应数据配置如下：

```
"backend": [
    {
      "url_pattern": "/",
      "encoding": "json",
      "extra_config": {
        "github.com/devopsfaith/krakend-httpcache": {}
      },
      "sd": "dns"
    }
  ]
```

注意：开启数据缓存之后，数据会被缓存在KrakenD的内存中，这将极大地消耗网关的内存，也会影响其性能，所以要谨慎使用。


#### Traffic Shadowing ####

KrakenD可以配置中间件让其忽略某个Backend返回的流量，这个功能一般使用场景为：某个应用上线了新版本，向用实际客户端流量对其进行测试，但是又不想将新版本的流量返回给客户端。

示例：
```
{
    "endpoint": "/user/{id}",
    "timeout": "150ms",
    "backend": [
        {
            "host": [ "http://my.api.com" ],
            "url_pattern": "/v1/user/{id}"
        },
        {
            "host": [ "http://my.api.com" ],
            "url_pattern": "/v2/user/{id}"
            "extra_config": {
                "github.com/devopsfaith/krakend/proxy": {
                    "shadow": true
                }
            }
        }
    ]
},
```

#### Flatmap ####

KrakenD可以对数组进行一些操作，包括修改Key和删除某个数组。当Flatmap开启之后，`group`和`target`仍能正常工作，但是`whitelist`、`blacklist`和`mapping`将被忽略。

示例：
```
  "extra_config": {
        "github.com/devopsfaith/krakend/proxy": {
            "flatmap_filter": [
                {
                    "type": "move",
                    "args": ["schools.42.students", "alumni"]
                },
                {
                    "type": "del",
                    "args": ["schools"]
                },
                {
                    "type": "del",
                    "args": ["alumni.*.password"]
                },
                {
                    "type": "move",
                    "args": ["alumni.*.PK_ID", "alumni.*.id"]
                }
            ]
        }
    }
```

官方有一个非常详细的例子来说明其具体工作方式，[点击我查看](https://www.krakend.io/docs/backends/flatmap/)。


## 认证 ##

### JWT Validation ###

KrakenD支持为Endpoint配置JWT认证，配置示例：
```
"endpoint": "/foo"
"extra_config": {
    "github.com/devopsfaith/krakend-jose/validator": {
        "alg": "RS256",
        "jwk-url": "https://url/to/jwks.json",
        "cache": true,
        "audience": [
            "audience1"
        ],
        "roles_key": "department",
        "roles": [
            "sales",
            "development"
        ],
        "issuer": "http://my.api.com",
        "cookie_key": "TOKEN",
        "disable_jwk_security": true,
        "jwk_fingerprints": [
            "S3Jha2VuRCBpcyB0aGUgYmVzdCBnYXRld2F5LCBhbmQgeW91IGtub3cgaXQ=="
        ],
        "cipher_suites": [
            10, 47, 53
        ]
    }
}
```

- **alg**: recognized string. The hashing algorithm used by the issuer. Usually RS256.
- **jwk-url**: string. The URL to the JWK endpoint with the public keys used to verify the authenticity and integrity of the token.
- **cache**: boolean. Set this value to true to store the JWK public key in-memory for the next 15 minutes and avoid hammering the key server, recommended for performance. The cache can store up to 100 different public keys simultaneously.
- **audience**: list. Set when you want to reject tokens that do not contain an audience of the list.
- **roles_key**: When passing roles, the key name inside the JWT payload specifying the role of the user.
- **roles**: list. When set, the JWT token not having at least one of the listed roles are rejected.
- **issuer**: string. When set, tokens not matching the issuer are rejected.
- **cookie_key**: string. Add the key name of the cookie containing the token when is not passed in the headers
- **disable_jwk_security**: boolean. When true, disables security of the JWK client and allows insecure connections (plain HTTP) to download the keys.
- **jwk_fingerprints**: string list. A list of fingerprints (the unique identifier of the certificate) for certificate pinning and avoid man in the middle attacks. Add fingerprints in base64 format.
- **cipher_suites**: integers list. Override the default cipher suites. Unless you have a legacy JWK, you don’t need to add this value.

支持的Hash算法（alg）:
```
EdDSA: EdDSA
HS256: HS256 - HMAC using SHA-256
HS384: HS384 - HMAC using SHA-384
HS512: HS512 - HMAC using SHA-512
RS256: RS256 - RSASSA-PKCS-v1.5 using SHA-256
RS384: RS384 - RSASSA-PKCS-v1.5 using SHA-384
RS512: RS512 - RSASSA-PKCS-v1.5 using SHA-512
ES256: ES256 - ECDSA using P-256 and SHA-256
ES384: ES384 - ECDSA using P-384 and SHA-384
ES512: ES512 - ECDSA using P-521 and SHA-512
PS256: PS256 - RSASSA-PSS using SHA256 and MGF1-SHA256
PS384: PS384 - RSASSA-PSS using SHA384 and MGF1-SHA384
PS512: PS512 - RSASSA-PSS using SHA512 and MGF1-SHA512
```
支持的密钥算法套件(cipher_suites):
```
5: TLS_RSA_WITH_RC4_128_SHA
10: TLS_RSA_WITH_3DES_EDE_CBC_SHA
47: TLS_RSA_WITH_AES_128_CBC_SHA
53: TLS_RSA_WITH_AES_256_CBC_SHA
60: TLS_RSA_WITH_AES_128_CBC_SHA256
156: TLS_RSA_WITH_AES_128_GCM_SHA256
157: TLS_RSA_WITH_AES_256_GCM_SHA384
49159: TLS_ECDHE_ECDSA_WITH_RC4_128_SHA
49161: TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
49162: TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
49169: TLS_ECDHE_RSA_WITH_RC4_128_SHA
49170: TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
49171: TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
49172: TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
49187: TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
49191: TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
```
默认密钥算法套件为：
```
49199: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
49195: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
49200: TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
49196: TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
52392: TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
52393: TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
```

### JWT Signing ###

配置示例：
```
"endpoints": [
    {
      "endpoint": "/token",
      "method": "POST",
      "extra_config": {
        "github.com/devopsfaith/krakend-jose/signer": {
          "alg": "HS256",
          "jwk-url": "http://backend/jwk/symmetric.json",
          "keys-to-sign": [
            "access_token",
            "refresh_token"
          ],
          "kid": "sim2",
          "cipher_suites": [
            5,
            10
          ],
          "jwk_fingerprints": [
            "S3Jha2VuRCBpcyB0aGUgYmVzdCBnYXRld2F5LCBhbmQgeW91IGtub3cgaXQ=="
          ],
          "full": true,
          "disable_jwk_security": true
        }
      }
    }
  ]
```

配置必填项：

- **alg**: recognized string. The hashing algorithm used by the issuer. Usually RS256.
- **jwk-url**: string. The URL to the JWK endpoint with the private keys used to sign the token.
- **kid**: string. The key ID member purpose is to match a specific key, as the jwk-url might contain several keys.
- **keys-to-sign**: string list. List of all the specific keys that need signing.

配置可选项：

- **full**: boolean. Use JSON format instead of the compact form JWT is giving.
- **disable_jwk_security**: boolean. When true, disables security of the JWK client and allows insecure connections (plain HTTP) to download the keys.
- **jwk_fingerprints**: string list. A list of fingerprints (the unique identifier of the certificate) for certificate pinning and avoid man in the middle attacks. Add fingerprints in base64 format.
- **cipher_suites**: integers list. Override the default cipher suites. Unless you have a legacy JWK, you don’t need to add this value.


### OAuth 2.0支持 ###

某些情况下，Backend需要认证才能访问，KrakenD支持OAuth 2.0认证。配置示例：

```
{
    "endpoint": "/endpoint",
    "backend": [
        {
            "url_pattern": "/backend",
            "extra_config": {
                "github.com/devopsfaith/krakend-oauth2-clientcredentials": {
                    "client_id": "YOUR-CLIENT-ID",
                    "client_secret": "YOUR-CLIENT-SECRET",
                    "token_url": "https://custom.auth0.tld/token_endpoint",
                    "endpoint_params": {
                        "client_id": ["YOUR-CLIENT-ID"],
                        "client_secret": ["YOUR-CLIENT-SECRET"],
                        "audience": ["YOUR-AUDIENCE"]
                    }
                },
                "github.com/devopsfaith/krakend-martian": {
                    "fifo.Group": {
                        "scope": ["request", "response"],
                        "aggregateErrors": false,
                        "modifiers": [
                            {
                                "header.Modifier": {
                                    "scope": ["request"],
                                    "name" : "Accept",
                                    "value" : "application/json"
                                }
                            }
                        ]
                    }
                }
            }
        }
    ]
}
```

- **client_id** string: The Client ID provided to the Auth server
- **client_secret** string: The secret string provided to the Auth server
- **token_url** string: The endpoint URL where the negotiation of the token happens
- **scopes**: string,optional A comma separated list of scopes needed, e.g.: scopeA,scopeB
- **endpoint_params** list,optional: Any additional parameters that you want to include in the payload when requesting the token. For instance, it is frequent to add the audience request parameter that denotes the target API for which the token should be issued.


## 服务发现 ##

### 静态 ###

在Backend例配置`"sd": "static"`即采用静态方式，配置如下：

```
"backend": [
	{
		"url_pattern": "/some-url",
		"sd": "static",
		"host": [
			"http://my-service-01.api.com:9000",
			"http://my-service-02.api.com:9000"
		]
	}
]
```

### DNS SRV ####

配置如下：

```
"backend": [
        {
          "url_pattern": "/foo",
          "sd": "dns",
          "host": [
            "api-catalog.service.consul.srv"
          ],
          "disable_host_sanitize": true
        }
      ]
```

### ETCD ###

首选在Root级配置中开启ETCD：
```
{
  "version": 2,
  "extra_config": {
    "github_com/devopsfaith/krakend-etcd": {
      "machines": [
        "https://192.168.1.100:4001",
        "https://192.168.1.101:4001"
      ],
      "dial_timeout": "5s",
      "dial_keepalive": "30s",
      "header_timeout": "1s",
      "cert": "/path/to/cert",
      "key": "/path/to/cert-private-key",
      "cacert": "/path/to/CA-cert"
    }
  },
```
其中只有`machines`是必须配置的，其余均为可选项。

开启了ETCD之后，可以在Backend中配置ETCD服务发现：
```
"backend": [
        {
          "url_pattern": "/foo",
          "sd": "etcd",
          "host": [
            "api-catalog"
          ],
          "disable_host_sanitize": true
        }
      ]
```

### Eureka ###

`Eureka`服务发现中间件并没有内置进KrakenD-CE，[参考文档](https://www.krakend.io/docs/service-discovery/eureka/)


## 限速与节流 ##

### 超时 ###

超时时间单位：ns、us/µs、ms、s、m、h。

**连接超时**

`timeout`可以为全局参数，也可以为某个Endpoint指定超时时间。当达到超时时间后，连接会被中断，客户端会收到`500 Internal Server Error`。配置示例：
```
{
  "version": 2,
  "timeout": "2000ms",
  "endpoints": [
    {
      "endpoint": "/splash",
      "method": "GET",
      "timeout": "1s"
      ...
    }
  ]
}
```
**HTTP请求超时**

- HTTP Read Timeout：指定读取整个HTTP请求的超时时间。
```
{
    "version": 2,
    "read_timeout": "1s"
}
```
- HTTP Write Timeout:指定组织响应体的超时时间。
```
{
    "version": 2,
    "write_timeout": "0s"  (no timeout)
}
```
- HTTP Idle Timeout:开启keep-alives后，指定等待下一个请求的超时时间。
```
{
    "version": 2,
    "idle_timeout": "0s"  (no timeout)
}
```
- HTTP Read Header Timeout:指定读取请求头部的超时时间。
```
{
    "version": 2,
    "read_header_timeout": "10ms"
}
```

### 控制闲置连接 ###

配置了`max_idle_connections`之后，当闲置连接达到所配数值之后，KrakenD会关闭那些处于keep-alive状态的连接，如果没有配置，KrakenD默认的闲置连接是250。


## Logging和Metrics ##

### Logging ###

配置示例：

```
{
  "version": 2,
  "extra_config": {
    "github_com/devopsfaith/krakend-gologging": {
      "level": "INFO",
      "prefix": "[KRAKEND]",
      "syslog": true,
      "stdout": true,
      "format": "custom",
      "custom_format": "%{message}"
    }
  }
```
日志等级：DEBUG、INFO、WARNING、ERROR和CRITICAL。

**输出日志至Graylog Cluster**

配置示例：

```
"extra_config": {
  "github_com/devopsfaith/krakend-gelf": {
    "address": "myGraylogInstance:12201",
    "enable_tcp": false
  }
  "github_com/devopsfaith/krakend-gologging": {
      "level": "INFO",
      "prefix": "[KRAKEND]",
      "syslog": false,
      "stdout": true
  }
}
```

- address: Graylog cluster的地址
- enable_tcp: 配置为false时采用UDP报文。

注意：同时需要开启`krakend-gologging`。

**输出Logstash格式日志**

配置示例：
```
"extra_config": {
  "github_com/devopsfaith/krakend-logstash": {
    "enabled": true
  }
  "github_com/devopsfaith/krakend-gologging": {
      "level": "INFO",
      "prefix": "[KRAKEND]",
      "syslog": false,
      "stdout": true,
      "format": "logstash"
  }
}
```

注意：同时需要开启`krakend-gologging`。


### Metrics ###

开启Metrics之后，KrakenD会新增一个Endpoint，监听一个不同的端口，URI为`/__stats/`。配置示例：

```
{
  "version": 2,
  "extra_config": {
    "github_com/devopsfaith/krakend-metrics": {
      "collection_time": "60s",
      "proxy_disabled": false,
      "router_disabled": false,
      "backend_disabled": false,
      "endpoint_disabled": false,
      "listen_address": ":8090"
    },
    ...
  }
```

- **collection_time**: The time window to collect metrics. Defaults to 60 seconds.
- **proxy_disabled**: Skip any metrics happening in the proxy layer (traffic against your backends)
- **router_disabled**: Skip any metrics happening in the router layer (activity in KrakenD endpoints)
- **backend_disabled**: Skip any metrics happening in the backend layer.
- **endpoint_disabled**: Do not publish the /__stats/ endpoint. Metrics won’t be accessible via the endpoint but still collected.
- **listen_address**: Change the listening address where the metrics endpoint is exposed. It defaults to :8090.