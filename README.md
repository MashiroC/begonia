# Begonia RPC

## 快速开始

### 概述

一个轻量级、API友好的RPC框架。目前仅支持Golang语言，Java、Node.js的支持正在计划中。

### 功能

- [x] client端 远程调用
- [x] server端 注册服务
- [x] 使用代码生成加速
- [x] 服务中心
- [ ] 日志库、日志中心
- [ ] 全链路追踪
- [ ] 连接池
- [ ] 心跳机制、机器监控
- [ ] 配置中心、多节点同步配置
- [ ] 服务管理、自动扩容

### 证书

`begonia` 的源码允许用户在遵循 [MIT 开源证书](https://github.com/MashiroC/begonia/blob/master/LICENSE) 规则的前提下使用。

### 安装

```bash
go get -u github.com/MashiroC/begonia
go get -u github.com/MashiroC/begonia/cmd/center
go get -u github.com/MashiroC/begonia/cmd/begonia
```

此时会下载begonia-rpc框架和begonia的命令行工具。推荐使用`go mod`。

### 快速开始

详细的使用文档请点击 [Begonia 使用文档](https://github.com/MashiroC/begonia/blob/master/doc/Begonia_Using_Doc) 。

框架中的架构、设计细节的设计文档正在编写中。

#### 服务模式

(非必须)

如果需要服务中心来调度服务、提供远程的服务注册和内网穿透，可以使用服务中心模式。

在安装后，执行下面的命令来启动服务中心的进程。

```bash
$ bgacenter start
```

推荐使用脚本来将服务中心注册为一个系统服务。(脚本在写了在写了)

如果没有上述需求，可以像gRPC一样直接从客户端连接服务。

#### 服务端

```go
// server.go
package main

import (
	"errors"
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"time"
)

func main() {
    // 一般情况下，addr是服务中心的地址。
    s := begonia.NewServer(option.Addr(":12306")) 
    
	echoService := &EchoService{}

    // 会通过反射的方式把EchoService下面所有公开的方法注册到Echo服务上。
	s.Register("Echo", echoService)

	s.Wait()
}

type EchoService struct {}

// SayHello 函数的参数和返回值会被反射解析，注册为一个远程函数。
// 注册的函数没有特定的格式和写法。
func (*EchoService) SayHello(name string) string {
	fmt.Println("sayHello")
	return "😀Hello " + name
}

```

以上这份代码的调用和解析都是使用反射来实现的，如果想要更快的速度，更高的并发，可以使用代码生成模式。代码生成不需要你提供额外的描述文档，也不需要特定的格式，仅需要切换到你项目的跟目录，执行以下命令：

```bash
$ begonia -s -r ./
```

之前下载的命令行工具会自动扫描调用了`begonia.Register`函数的结构体。并在你结构体的同级目录生成`${structName}.begonia.go`文件，并且在初始化时使用代码生成相关的API，中间任何阶段不使用反射。

```bash
└─server
    ├─server.go
    └─EchoService.begonia.go
```

#### 客户端

```go
package main

import (
	"fmt"
	"github.com/MashiroC/begonia"
	"github.com/MashiroC/begonia/app/option"
	"sync"
	"time"
)

const (
	mode = "center"
)

func main() {
    c := begonia.NewClient(option.Addr(":12306"))

    // 获取一个服务
    s, err := c.Service("Echo")
	if err != nil {
		panic(err)
	}

    // 获取一个远程函数的同步调用
	testFun, err := s.FuncSync("SayHello")
	if err != nil {
		panic(err)
	}
    
    // 获得一个远程函数的异步调用
    testFunAsync, err := s.FuncAsync("SayHello")
    if err != nil {
		panic(err)
	}
    
    // 调用！
	testFunAsync(func(res interface{}, err error) {
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}, "kieran")

	res, err = testFun("kieran")
	if err != nil {
		panic(err)
	}
    fmt.Println(res)
}

```

以上这个是使用反射调用的例子，在服务端可以使用代码生成直接生成调用文件。

```bash
$ begonia -s -c -r ./
```

执行后在服务的目录下会生成一个`call`目录，目录包含所有服务的调用代码。其中`call.begonia.go`是begonia客户端的初始化代码，可以修改替换。`entity.begonia.go`是注册服务上的结构体的声明。

```bash
└─server
	├─call
		├─cli.begonia.go
		├─entity.begonia.go
		└─EchoService.begonia.go
    ├─server.go
    └─EchoService.begonia.go
```

客户端可以将这段代码复制到自己的项目里直接调用：

```go
package main

import (
	"fmt"
	"github.com/MashiroC/begonia/example/service/call"
)

func main() {
	res, err := call.SayHello("kieran")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```