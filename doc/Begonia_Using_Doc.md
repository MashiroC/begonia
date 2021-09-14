# Begonia RPC 使用文档

## 前言

begonia是一款轻量级、快速并且简单的RPC框架，目前仅提供了远程调用和服务注册的能力，接下来begonia会提供一些更高级的功能。

### 重要概念

#### 一切皆RPC

在begonia内，无论是你编写的RPC请求、还是框架内的服务注册、获取服务信息、抑或是后续的同步日志、同步配置等功能。皆为一个RPC请求。

#### 节点

beognia将节点分为两类：客户端(client)、服务端(server)以及混合节点(mix)。

客户端即只会进行远程调用，而不会注册服务的节点。

服务端即只会注册服务，而不会进行远程调用的节点。

> 其实服务端在注册服务时会调用一个服务REGISTER.Register。所以实际上服务端是拥有远程调用的功能的，但是在sdk的api上隐藏掉了，混合节点只不过是开放了远程调用相关的api。

### 集群模式

#### Center Cluster

当你使用快速开始中的例子时，便是使用了`Center Cluster`模式。该模式会启动一个服务中心，其他节点都会向服务中心发起一条TCP连接。这个center进程就是一个服务中心。

center是一个特殊的server节点，这个节点会注册一系列服务，例如服务注册服务"REGISTER"，日志服务"LOG"，监控服务"WATCH"等(除了注册，其他的在写了在写了)。这些服务用于在你的集群上提供高级功能。

并且center节点会开启"代理职责"用于将一个调用转发到正确的节点。

#### P2P

类似gRPC，每个服务都进行监听，不在远程服务中心注册。客户端需要针对不同的服务建立不同的连接。

如果想使用这种服务模式，仅需要在server的初始化阶段传入一个配置`option.P2P()`即可

#### P2P Cluster

经典的服务发现模式，每个服务都进行监听，并且将服务注册在etcd / zookeeper上。LOG、WATCH等服务启动会单独启动监听。

(正在开发中)

## 初始化

begonia通过`begonia.NewClient()`来创建一个客户端实例。同样的，通过`begonia.NewServer()`创建一个服务端实例。

这两个函数的参数都是`...option.WriteFunc`。诸如监听(或连接)地址、节点模式这些信息通过option来传递给实例。而例如接收一个新包时的超时时间等，目前全部写死在了`config`包中，接下来会提供一个配置这些细则的能力。

目前option中仅支持下列两个配置，配置通过option包中的函数传入

| 函数名              | 说明                                                         |
| ------------------- | ------------------------------------------------------------ |
| option.Addr(string) | 地址。对于client节点是服务中心地址或服务地址。对于service节点根据配置来连接到某地址或监听某地址。 |
| option.P2P()        | 点对点连接。此配置只有服务端调用会有效果，服务端会在本地根据option.Addr传入的地址进行监听 |

推荐使用默认的服务中心模式，如果使用默认模式，请在开始服务端和客户端的开发前启动服务中心。

```bash
$ bgacenter start
```

如果仅拥有一个server节点，可以使用p2p模式，请在服务端启动时传入`option.P2P()`，此时不需要启动服务中心。

## 服务端

### 注册服务

在begonia里，服务注册非常简单。你只需要调用这样一个函数就可以将一个服务注册到服务中心或本地。

```go
s.Register("Echo", &echoService{})
```

这个函数的函数签名如下：

```go
Register(name string, service interface{}, registerFunc ...string)
```

这个函数有多个参数：

第一个是注册在服务中心的服务名，客户端会根据这个名字来查找服务。

第二个是一个结构体的指针，例如我们传递了一个`echoService`的指针，那么这个`echoService`下的所有公开函数都会被注册为远程函数。

最后是一个可变列表，如果你需要指定这个结构体下哪些函数被注册，当你忽略这个变量时，解析器会将结构体下所有函数注册。只要你传入了一个函数名，那么解析器就只会解析你传入的函数名了。

### 远程函数

begonia中对于远程函数没有强制的格式，只要是结构体下的公开函数都可以被注册到服务上。例如以下这些都可以被解析器解析后在服务中心注册。

```go
func (*echoService) SayHello(name string) string {
return "Hello " + name
}

func (*EchoService) SayHelloWithContext(ctx context.Context, name string) string {
return "Hello ctx " + name
}

func (*EchoService) Add(i1, i2 int) (res int, err error) {
res = i1 + i2
return
}

func (*EchoService) NULL() {}

func (EchoService) Hello() (err error) {
fmt.Println("Hello")
}
```

当一个函数为`context.Context`时，解析器会忽略它。在这个函数被调用时框架会自动生成一个`context.Context`，可以从这个`ctx`中获得本次函数调用的一些信息，例如这次请求的唯一ID、本次请求的连接的ID。可以通过下列示例代码获取这些信息：

```go
func (*EchoService) SayHelloWithContext(ctx context.Context, name string) string {
v:=ctx.Value("info")
info:=v.(map[string]string)
fmt.Println(info)
fmt.Println(info["reqID"])
fmt.Println(info["connID"])
return "Hello ctx " + name
}
```

当一个函数的最后一个返回值为`error`时，解析器会自动解析到客户端的`error`上，这一部分请转向继续阅读客户端部分。

目前begonia的远程函数可以支持的类型有：

- `int` , `[]int`
- `int8` ~ `int64` , `[]int8` ~ `[]int64`
- `float32` , `float64` , `[]float32` , `[]float64`
- `bool` , `[]bool`
- `string` , `[]string`
- `struct` , `[]struct`
- `[]byte`
- `map ` ( 不支持`map[string]interface{}` )
- `error` ( `error`类型推荐只放在函数最后一个返回值上 )

### 代码生成

begonia服务端的代码生成并不会改变注册的服务的写法，也不需要你编写额外的代码、描述文件。所以你需要做的仅仅是在编译前运行一下代码生成工具而已:

```bash
$ begonia -s -r ./
generate server1 example\server\server.go_EchoService ...
server1 code ok!
complete, total: 1.3614423s
```

命令中`-s`代表生成服务端代码，`-r`代表移除旧的生成代码。`./`代表扫描当前目录。该命令行工具会自动扫描目录下所有调用了`Server.Register()`注册的结构体，自动为它们添加代码替代反射来进行加速。调用后会在注册结构体的文件同级目录生成文件`${structName}.begonia.go`。

## 客户端

### 获取服务

如果你不希望使用代码生成，可以通过下列代码获得一个远程函数并调用：

```go
s, _ := c.Service("Echo")

SayHello, _ := s.FuncSync("SayHello")

res, err := SayHello("kieran")
fmt.Println(res, err)

SayHelloAsync , _ := s.FuncAsync("SayHello")

SayHelloAsync(func(res interface{}, err error) {
fmt.Println(res, err)
}, "kieran")
```

如上示例，你可以获得一个同步函数也可以获得一个异步函数。下面是这两种方式的函数签名：

```go
// RemoteFunSync 同步远程函数
type RemoteFunSync func(params ...interface{}) (result interface{}, err error)

// RemoteFunAsync 异步远程函数
type RemoteFunAsync func(callback AsyncCallback, params ...interface{})

// AsyncCallback 异步回调
type AsyncCallback = func(interface{}, error)
```

可以看出两种方式除了同步和异步之外没有区别，在文档里只会详细讲述同步函数。

### 调用函数

在反射模式下(即不使用代码生成)，调用`s.FuncSync()`或`s.FuncAsync()`后会生成一个匿名函数，该函数的参数是`interface{}`的可变列表，返回值是`interface{}`和`error`。

参数不需要指定参数名对应的值，但是需要和服务端的参数顺序相同。

例如一个远程函数是`(name string, age int)`。在这里你填的肯定不能是`(age,name)`

(我总觉得这里应该不会有人参数顺序写错吧...)

#### error

这个函数调用后会将参数序列化后传输到服务端并等待响应。如果远程函数最后一个返回值为`error`，该返回值会被解析到客户端返回的`error`上。该`error`的解析拥有一个优先级，如果同时出现了多个错误，会根据优先级赋值到客户端的`error`上：

```
框架报错 => 远程函数的panic => 远程函数的error返回值
```

接下来的文档将会专注于其他返回值，会忽略`error`，不再对其做特别描述。

#### interface{}

远程函数中非`error`的部分会被解析到客户端的`interface{}`上，本地的匿名函数返回值的类型会根据下表做解析：

| 远程函数                 | 匿名函数            |
| ------------------------ | ------------------- |
| int , int8 ~ int64       | int                 |
| []int , []int8 ~ []int64 | []int               |
| bool , [] bool           | bool , []bool       |
| float32 , float64        | float32 , float64   |
| []float32 , []float64    | []float32 , float64 |
| []byte                   | []byte              |
| string                   | string              |
| map                      | map                 |
| struct                   | struct              |

除此之外，根据返回值的数量，匿名函数的`interface{}`也会做相应变化

- 无返回值，例如：

  ```go
  func (*EchoService) Hello()
  ```

  这种函数的值会被解析为`bool`，如果函数调用成功则返回`true`。

  ```go
  Hello, _ := s.FuncSync("Hello")
  res, _ := SayHello()
  fmt.Println(res.(bool))
  ```

- 一个返回值，例如：

  ```go
  func (*EchoService) SayHello(name string) string
  ```

  这种函数的值会直接根据返回值类型，进行解析到`interface{}`上。

  ```go
  SayHello, _ := s.FuncSync("SayHello")
  res, _ := SayHello("kieran")
  fmt.Println(res.(string))
  ```

- 多返回值，例如：

  ```go
  func (*EchoService) Mod(i1, i2 int) (res1 int, res2 int)
  ```

  这种函数的返回值会被解析为`[]interface{}`。

  ```go
  mod, _ := s.FuncSync("Mod")
  res, _ := mod(5,2)
  resArr := res.([]interface{})
  fmt.Println(resArr[0].(int))
  fmt.Println(resArr[1].(int))
  ```

### 代码生成

相对来说，使用反射API的远程调用还是较为繁琐，你可以在服务端可以通过代码生成工具直接生成客户端调用相关代码，客户端可以将这一份代码复制到本地来直接调用远程服务。

在生成服务端代码的命令中加入一个`-c`命令生成客户端调用代码：

```bash
$ begonia -s -c -r ./
generate server1 example\server\server.go_EchoService ...
server1 code ok!
client call ok!
complete, total: 1.6967934s

$ tree ./
└─call
	├─cli.begonia.go
	├─entity.begonia.go
	└─EchoService.begonia.go
├─server1.go
└─EchoService.begonia.go
```

命令执行完成后，除了服务端相关代码，还会生成一个`call`目录，目录下有多个文件：

- cli.begonia.go

  客户端实例的初始化相关代码

- entity.begonia.go

  服务中使用的结构体的声明

- ${structName}.begonia.go

  服务对应的结构体生成的代码

客户端可以将这个包复制到自己的项目里直接调用：

```go
package main

import (
  "fmt"
  "github.com/MashiroC/begonia/example/server/call"
)

func main() {
  res, err := call.SayHello("kieran")
  if err != nil {
    panic(err)
  }
  fmt.Println(res)
}
```