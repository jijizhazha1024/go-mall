# go-mall

## 中间件使用：
### 必须设定`WithClientMiddleware,WrapperAuthMiddleware`两个中间件
[中间件](common/middleware)
使用示例：
api文件
```api
syntax = "v1"

type Request {}

type Response {
	Message string `json:"message"`
}

@server (
	middleware: WithClientMiddleware,WrapperAuthMiddleware // 必须按照先后顺序。
	prefix:     /douyin/test
)
service test-api {
	@handler TestHandler
	get / (Request) returns (Response)
}


```
servicecontext.go文件
```go
type ServiceContext struct {
    Config                config.Config
    WrapperAuthMiddleware rest.Middleware
    WithClientMiddleware  rest.Middleware
    }

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:                c,
        WrapperAuthMiddleware: middleware.WrapperAuthMiddleware(c.AuthsRpc), // # 需要指定认证rpc地址
        WithClientMiddleware:  middleware.WithClientMiddleware,
    }
}


```

### 从ctx中获取用户id和客户端IP
> 注意只能在api层获取到客户端IP和用户id
```
clientIP := r.Context().Value(biz.ClientIPKey).(string)
userId := r.Context().Value(biz.UserIdKey).(string)
```
## 运行脚本使用

项目目录下的run.go为服务运行脚本，使用如下：

```shell
go run run.go --services service1,service2
```

example1:

```shell
go run run.go --services user,goods
```

example2:

不进行指定services默认启动所有服务

```shell
go run run.go
```

