# Begonia 设计文档

## 概述

Begonia 是一款轻量级的 RPC 框架，目标是提供开箱即用的 RPC 与微服务框架，告别繁杂的配置与部署。Begonia 正在不断开发中，接下来会逐渐提供日志、链路追踪、服务监测、报警、自动扩缩容、服务限流、熔断等能力。

在 Begonia 中，一切皆 RPC ，一切皆服务。框架中一切的能力借由注册在服务中心的服务提供。甚至用于注册的服务中心也是一个服务。用户注册服务、日志同步等功能或模块，全部注册在服务中心上，用户通过一个 RPC 调用即可。

### 节点

Begonia 使用经典的 C (`Client`) / S (`Server`) 架构，我们称一个引用了 Begonia sdk 的进程为一个节点，节点分为 `Client` 节点和 `Server` 节点两种。我们建议使用容器来运行 Begonia 的节点，一个由多个 Begonia 节点组成的集群如下图所示：

![Begonia 集群网络拓扑](https://github.com/MashiroC/begonia/blob/master/doc/pic/begonia%20%E9%9B%86%E7%BE%A4%E7%BD%91%E7%BB%9C%E6%8B%93%E6%89%91.png)

这里除掉服务中心外一共有 6 个节点，这些节点都会向服务中心发起一条 TCP 连接，RPC 的请求和响应帧全部通过这条连接来传输。服务中心在这里起到的是服务发现、内网穿透、负载均衡等作用。

`Client` 节点不应注册任何服务，只调用 `Server` 节点的服务即可。

`Server` 节点应当挂载一些中台服务，供 `Client` 节点调用。服务中心同样是一个 `Server` 节点，它上面也注册了一些服务。但该节点拥有一个特殊的能力：代理职责。详细内容会在下面的部分细说。

### 请求流程

如果我们现在拥有一个 `Server` 节点，注册一个`Echo`服务，服务中有一个函数`SayHello()`，并且拥有一个 `Client` 节点，请求调用`Echo.SayHello`，流程如下：

![Begonia 流程图](https://github.com/MashiroC/begonia/blob/master/doc/pic/Begonia%20%E6%B5%81%E7%A8%8B%E5%9B%BE.png)

服务中心、客户端、服务端初始化之后便开始监听，其中服务中心会在开始监听前注册一个服务`REGISTER`用来提供服务发现能力。服务端注册服务通过一个 RPC 调用`Register.Register()`来将服务的相关信息注册到注册中心。客户端通过一个 RPC 调用`Register.ServiceInfo()`来获取服务的相关信息，最后生成匿名函数供用户调用。

当客户端获得服务的信息后，每次不需要重复黄色部分，只需要蓝色的调用部分就可以了。

### 架构

Begonia 从代码层面抽象成了三层，这三层分别提供不同的能力。

![Begonia 代码架构](https://github.com/MashiroC/begonia/blob/master/doc/pic/Begonia%20%E6%9E%B6%E6%9E%84.png)

**App (Application) 层**

这一层主要提供的是序列化、反序列化能力和正确的函数调用能力，用户调用的API就来自于这里。

**Logic 层**

`Logic` 层收到 `App` 层的包已经经过了序列化，所以在这里仅需要将其封装为后发送给 `Dispatch` 层，然后注册该请求的回调，等待响应到来。

并且 `Logic` 层会在 `Dispatch` 层中注册事件，当有新的数据包到来时，`Logic` 会通过数据包的序列号来查找相对应的回调。

**Dispatch 层**

`Dispatch` 层主要负责的是通讯，对端口的监听、或连接至对应的地址。

## starter

starter的作用主要是为了引导 `App` 、`Logic` 和 `Dispatch` 层初始化。

在Begonia中，初始化需要调用下面两个函数：

```go
// 客户端的初始化
begonia.NewClient()

// 服务端的初始化
begonia.NewServer()
```

查看这两个函数的函数签名，它们的入参是相同的：

```go
func NewClient(optionFunc ...option.WriteFunc) (cli Client)

func NewServer(optionFunc ...option.WriteFunc) (s Server)
```

它们的入参都是一些初始化选项，比如`option.Addr()`或`option.P2P()`。

这些选项生成的`option.WriteFunc`是一个函数，它的函数签名如下：

```go
type WriteFunc func(optionMap map[string]interface{})
```

这些函数的作用是往一个 `map` 中写入数据。

在上述两个初始化函数中，会首先得到一份默认选项，然后遍历并运行传入的`option.WriteFunc`，最后得到一份选项的 `map` ，传递给相应的`Client`或`Server`的初始化函数。初始化函数会根据 `map` 中的选项来做相应的初始化。

并且，当你调用这两个函数时，这两个函数位于根目录 `begonia` 包下，其中 `Server` 和 `Client` 的接口声明只不过是 `app/client` 或 `app/server` 包下的接口的一个 alias 。

```go
import (
	"github.com/MashiroC/begonia/app/option"
	"github.com/MashiroC/begonia/app/server"
)

// Server 服务端的接口
type Server = server.Server

// New 初始化
func NewServer(optionFunc ...option.WriteFunc) (s Server) {
	optionMap := defaultServerOption()

	for _, f := range optionFunc {
		f(optionMap)
	}

	in := server.BootStart(optionMap)
	return in
}
```

### 代码生成

`App` 层是使用代码生成的实现还是使用反射的实现，这个选项的配置位于 `app` 包下， `app` 包下有一个全局变量`ServiceAppMode`，默认是反射模式，starter在引导的时候会去根据这个值来判断该创建哪种实现的实例。

```go
var ServiceAppMode = Reflect

type ServiceAppModeTyp int

const (
	invalid ServiceAppModeTyp = iota
	Ast
	Reflect
)
```

当你使用代码生成的方式去生成 `Client` 或 `Server` 的时候，生成的代码中会在 `init()` 函数的第一行修改这个选项。

## App层

`App` 层有两个实现，`reflect` 和 `ast`。分别是使用反射和代码生成实现的 API。在一般情况下， 使用反射即可实现 RPC 的功能，但反射的开销较大，用户可以使用代码生成来替换反射的实现，去除反射的开销。两个实现的接口完全相同，并且不需要编写的描述文件或代码。

客户端的调用和服务端的响应在这里会封装为`logic.Call`和`logic.CallResult`后会传递给 `Logic` 层。`logic.Call`结构体的签名如下：

```go
type Call struct {
	Service string // 调用的服务名
	Fun     string // 调用的函数名
	Param   []byte // 远程函数的入参(已序列化)
}
```

一个客户端的调用会在 `App` 层将入参序列化为 `[]byte` 后传递给 `Logic` 层。同样的，在 `Logic` 层传递给 `App` 层的 `logic.CallResult` 中的出参也会在 `App` 层反序列化。

```go
type CallResult struct {
	Result []byte // rpc调用的结果，远程函数的出参(已序列化)
	Err    error  // 错误
}
```

在`CallResult`中已经包含了`Err`参数，所以在 `Server` 端的 `App` 层在序列化时会忽略返回值中最后一个 `error` 值，并将该`error`在框架不报错或没有捕获到 `panic` 的情况下直接传递给`CallResult.Err`。

### Server

`Server` 的接口声明如下：

```go
type Server interface {

	// Register 注册服务
	Register(name string, service interface{}, registerFunc ...string)

	// Wait 阻塞等待
	Wait()
	
}
```

#### 服务注册

当服务注册时，解析器会解析该函数的参数和返回值类型，最终生成一个 `avro schema`。`avro` 是一个数据序列化系统，只需要提供一个模式 ( `schema` ) ，可以通过 `schema` 将数据序列化为二进制格式在网络中传输。详细可以使用搜索引擎搜索 `apache avro` 。

在 begonia 中，函数的入参和出参的类型会被解析为 `avro schema` ，通过该 `schema`，客户端可以检验传入数据，服务端可以正确获得远程函数的入参。

在注册函数时，注册服务需要以下信息：

```go
type Service struct {
	Name string    // 服务名
	Mode string    // 服务的序列化模式，avro / protobuf
	Funs []FunInfo // 服务的函数
}

type FunInfo struct {
	Name      string // 函数名
	InSchema  string // 函数入参的avro schema
	OutSchema string // 函数出参的avro schema
}
```

这里的序列化模式目前只有 `avro` 一种，在将来会考虑加入更多的序列化模式。

#### 函数解析

在 `reflect` 模式下，上述需要的信息中 `InSchema` 和 `OutSchema` 是通过解析器解析函数获得的。还是假设我们有一个服务，服务下有一个 `AddAge` 函数：

```go
func (*TestService) AddAge(p People, num int) People {
    p.Age += num
    return p
}

// s.Register("Test",&TestService{})
```

通过解析后，得到一个入参的 `schema` ：

```json
{
    "namespace": "begonia.fun.AddAge",
    "type": "record",
    "name": "In",
    "fields": [
         {"name": "F1", "type": 
          	{
    			"type": "record",
    			"name": "People",
    			"fields":[
        			{"name": "Name", "type": "string"},
        			{"name": "Age",  "type": "int"} 
    			]
			}
         },
         {"name": "F2", "type": "int"}
    ]
}]
```

其中函数的入参 `p` 和 `num`，根据它们在入参中的顺序，`name` 值为 `F1` 和 `F2`。如果有更多的参数，也依次为`F3`、`F4`...

出参也类似，只不过 `name` 值变为 `Out`。

当一个函数解析完，获得了两个 `schema` 之后，便会调用注册器来完成注册。

#### 注册器

注册器其实是对注册过程的封装，它是一个接口：

```go
import cRegister "github.com/MashiroC/begonia/core/register"

type Register interface {
	
	// Register 注册一个函数
	Register(name string, info []cRegister.FunInfo) (err error)

	// Get 获得一个已经注册的函数，如果没有查找到，会返回一个error
	Get(name string) (fs []cRegister.FunInfo, err error)

}
```

这个接口有两个实现，分别是本地注册器 ( `localRegister` ) 和远程注册器 ( `remoteRegister` )。对于服务中心或使用了P2P模式的服务，服务注册在本节点的 `REGISTER` 服务上，在初始化的时候需要传入一个 `REGISTER` 服务的指针。

```go
// localRegister 本地注册器
type localRegister struct {
	c *cRegister.CoreRegister
}

func NewLocalRegister(c *cRegister.CoreRegister) Register {
	return &localRegister{
		c: c,
	}
}
```

除了服务中心之外的节点则使用远程注册器，远程注册器会将当前节点的注册请求封装为 RPC 发送到服务中心，当该 RPC 成功后，即代表注册完成。

```go
// remoteRegister 远程注册器
type remoteRegister struct {
	lg *logic.Client
}

func NewRemoteRegister(lg *logic.Client) Register {
	return &remoteRegister{
		lg: lg,
	}
}
```

注册后，服务端会进一步封装，将 `schema` 封装为 `coding.Coder`。

#### coding.Coder

`coding.Coder` 的作用是编解码器，它是一个接口，目前只有一个实现：

```go
type Coder interface {
	Encode(data interface{}) ([]byte, error)
	Decode([]byte) (data interface{}, err error)
	DecodeIn([]byte, interface{}) error
}

// avroCoder Avro模式的coder
type avroCoder struct {
	Schema avro.Schema
}
```

`avroCoder` 会根据 `avro schema` 来将一个结构体序列化或反序列化。上述的两个 `schema` 就会封装为 `InCoder` 和 `OutCoder`，供服务端在收到调用请求时反序列化和返回结果时序列化。

#### 代码生成

当使用代码生成的手段为服务加速后，并不需要改变调用，`Server` 在初始化时会初始化 `astServer` 而非` rServer`，`astServer` 中执行函数不再使用反射的方式调用函数。

代码生成工具会为服务的结构体实现接口：

```go
// CodeGenService 代码生成实现的服务
type CodeGenService interface {

	// Do 调用服务
	Do(ctx context.Context, fun string, param []byte) (result []byte, err error)

	// FuncList 返回要注册的函数
	FuncList() []coreRegister.FunInfo
	
}
```

`App` 层在收到一个请求后，会直接调用服务的 `Do` 函数，`Do` 函数中是自动生成的序列化和反序列化与调用函数，下面是生成代码的一个示例：

```go
func (d *EchoService) Do(ctx context.Context, fun string, param []byte) (result []byte, err error) {
	switch fun {

	case "SayHello":
		var in _EchoServiceSayHelloIn
		err = _EchoServiceSayHelloInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1 := d.SayHello(
			in.F1,
		)
		if err != nil {
			return nil, err
		}
        
		var out _EchoServiceSayHelloOut
		out.F1 = res1

		res, err := _EchoServiceSayHelloOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	default:
		err = errors.New("rpc call error: fun not exist")
	}
	return
}
```

### Client

Client的接口声明如下：

```go
type Client interface {
	Service(name string) (s Service, err error)
	FunSync(serviceName, funName string) (rf RemoteFunSync, err error)
	FunAsync(serviceName, funName string) (rf RemoteFunAsync, err error)
	Wait()
	Close()
}
```

用户可以通过 `Service()` 函数来获取一个服务，再通过 `FunSync()` 或 `FunAsync()` 来获取一个远程函数。也可以直接在 `Client` 上调用这两个函数来获取远程函数。

#### 服务发现

当用户需要发现一个服务时，可以发起 RPC 调用 `REGISTER.ServiceInfo()` 来获取一个服务的信息，该调用会发送至服务中心查询服务相关信息，然后将服务下的函数信息返回。

获取信息这一步抽象为了通过注册器获得信息，通常情况下，使用远程注册器会发起 RPC 调用。而本地注册器则不会发起 RPC 调用。

```go
// Service 获取一个服务
func (r *rClient) Service(serviceName string) (s Service, err error) {
	
	fs, err := r.register.Get(serviceName)
	if err != nil {
		return
	}

	funs := make([]Fun, 0, len(fs))

	for _, f := range fs {
		inCoder, err := coding.NewAvro(f.InSchema)
		if err != nil {
			return nil, err
		}

		outCoder, err := coding.NewAvro(f.OutSchema)
		if err != nil {
			return nil, err
		}
		funs = append(funs, Fun{
			Name:     f.Name,
			InCoder:  inCoder,
			OutCoder: outCoder,
		})
	}

	if app.ServiceAppMode == app.Ast {
		s = r.newAstService(serviceName, r)
	} else {
		s = r.newService(serviceName, funs)
	}

	return
}
```

获得到信息后便会生成对应服务，继而获得匿名函数。

#### 匿名函数

获得的匿名函数的作用是将入参根据顺序生成一个 map 后，交给 `Coder` 去编码。编码完成便生成 `logic.Call` 交给Logic层做进一步处理。

TCP 在通讯中只能以异步的方法等待响应。在这里创建一个管道等待回调传回数据即可。

```go
rf = func(params ...interface{}) (result interface{}, err error) {
		ch := make(chan *logic.CallResult)

		b, err := f.InCoder.Encode(coding.ToAvroObj(params))
		if err != nil {
			return
		}

		r.c.lg.CallAsync(&logic.Call{
			Service: r.name,
			Fun:     name,
			Param:   b,
		}, func(res *logic.CallResult) {
			ch <- res
		})

		tmp := <-ch
		if tmp.Err != nil {
			err = tmp.Err
			return
		}

		// 对出参解码
		out, err := f.OutCoder.Decode(tmp.Result)

		result = reflects.ToInterfaces(out.(map[string]interface{}))
		return
	}
```

#### 代码生成

使用代码生成工具生成的代码中，会自动创建函数，如需使用可以直接将代码(或包)复制到客户端仓库。

生成的代码会实现编解码相关功能和封装。通过 `s.FuncSync` 获得的只不过是 `Logic` 层相关功能的一个转发。

`ast` 模式和 `reflect` 模式在这里的区别是从 `Client` 获取到的 `Service` 不同，一个是 `astService`，另一个是`rService`，里面对于获取到的匿名函数有不同的实现。

### 服务中心

服务中心是一个特殊的 `Server`。它在初始化的时候注册器为本地注册器 `localRegister`，在初始化完成，开始对端口进行监听前，会注册 `REGISTER` 服务。该服务用于其他 `Server` 节点注册服务。

#### REGISTER

`REGISTER` 本身也是一个代码生成的服务，它的代码如下：

```go
type CoreRegister struct {
	services *registerServiceStore
}

func (r *CoreRegister) Register(ctx context.Context, si Service) (err error) {
	v := ctx.Value("info")
	info := v.(map[string]string)
	connID := info["connID"]

	err = r.services.Add(connID, si.Name, si.Funs)
	return err
}

func (r *CoreRegister) ServiceInfo(serviceName string) (si Service, err error) {
	service, ok := r.services.Get(serviceName)
	if !ok {
		err = fmt.Errorf("server [%s] not found", serviceName)
		return
	}

	si.Funs = service.funs
	si.Name = serviceName
	si.Mode = "avro"
	return
}
```

它会通过 `ctx` 得到当前请求注册的`连接ID`，然后将相关信息存储在本地。

并且服务中心在初始化时会创建一个 `ProxyHandler` 注入 `Dispatch` 中。

#### ProxyHandler

`ProxyHandler` 为一个节点赋予代理职责，它的代码非常简短：

```go
type HandlerAction = func(connID, redirectConnID string, f frame.Frame)

type CheckFunc = func(connID string, f frame.Frame) (redirectConnID string, ok bool)

func NewHandler() *Handler {
	return &Handler{}
}

type Handler struct {
	Check CheckFunc

	handlerChains []HandlerAction
}

func (c *Handler) AddAction(action HandlerAction) {
	if c.handlerChains == nil {
		c.handlerChains = make([]HandlerAction, 0, 2)
	}

	c.handlerChains = append(c.handlerChains, action)
}

func (c *Handler) Action(connID, redirectConnID string, f frame.Frame) {
	if c.handlerChains == nil {
		return
	}

	for i := 0; i < len(c.handlerChains); i++ {
		c.handlerChains[i](connID, redirectConnID, f)
	}
}
```

其中的 `Check` 函数在创建的时候传入，`Check` 函数用来检查该请求是否应该被转发(或丢弃)。如果检查返回`true`，则会继续执行下一步。服务中心正是使用了 `ProxyHandler` 来为非 `REGISTER` 服务的请求转发到正确的连接。

## Logic层

对于 `App` 层，`Server` 节点仍然会调用注册相关 RPC 请求，但正常情况下 `Client` 不会收到 RPC 调用。在 `App` 层， 将两个节点区分开，但在 `Logic` 层的关系是 `Server` 包含 `Client`。

### Client

```go
// Client Client接口的实现结构体
type Client struct {
	Dp        dispatch.Dispatcher // dispatch层的接口，供logic层向下继续调用
	Callbacks *CallbackStore      // 回调仓库，可以在这里注册回调，调用回调
}
```

`Client` 中有两个参数：`Dp` 和 `Callbacks`。`Dp` 是 `Logic` 层持有的 `Dispatch` 层的指针，`Callbacks`是回调仓库。

#### 同步与异步请求

`App` 层的调用传递到 `Logic` 层中，会通过`Client.CallSync()`函数进行同步请求，通过`Client.CallAsync()`进行异步请求。在这里会首先生成一个唯一的 `请求ID`，目前该 ID 使用 `uuid` 实现。生成后创建用于 `Dispatch` 层通讯的 `frame.Frame` 。并向回调仓库添加响应到来时的回调，最后发送该请求到网络。

```go
type Callback = func(result *CallResult)

func (c *Client) CallAsync(call *Call, callback Callback) {

	reqID := ids.New()
    
	var f frame.Frame
	f = frame.NewRequest(reqID, call.Service, call.Fun, call.Param)

	c.Callbacks.AddCallback(context.TODO(), reqID, callback)

	if err := c.Dp.Send(f); err != nil {
		err = c.Callbacks.Callback(reqID, &CallResult{
			Result: nil,
			Err:    fmt.Errorf("logic call error: %w", err),
		})
		if err != nil {
			// TODO:println => errorln
			log.Println(err)
		}
	}

}
```

TCP 的本质是异步的，但是我们提供了同步的函数。同步函数本质上是异步函数的一个封装，通过管道来等待请求完成。

```go
func (c *Client) CallSync(call *Call) *CallResult {

	ch := make(chan *CallResult)
	defer close(ch)

	c.CallAsync(call, func(result *CallResult) {
		ch <- result
	})

	return <-ch
}
```

#### 回调仓库

`Logic` 层的核心是回调仓库 `CallbackStore`。

```go
// CallbackStore 回调仓库，拥有注册回调、回调的方法
type CallbackStore struct {
	chLock sync.RWMutex                // 锁
	ch     map[string]chan *CallResult // 存储的map
}

// AddCallback 添加一个回调
func (w *CallbackStore) AddCallback(ctx context.Context, reqID string, callback Callback)

// Callback 注册后调后，根据reqID来调用回调
func (w *CallbackStore) Callback(reqID string, cr *CallResult) (err error)

// goWait 这个需要开一个新协程 来等待结果或者超时
func (w *CallbackStore) goWait(reqID string, timeout, parent <-chan struct{}, cb Callback, ch chan *CallResult)
```

每当一个请求发出，就会在这里添加一个回调，等待返回值到达。返回值根据唯一的 `reqID` 来找到对应等待的回调，将数据发送回去。

在添加回调的时候，同时会新起一个协程用于超时，当一定时间内没有收到响应，便会认为丢包，返回超时。

### Server

在 `Logic` 层，`Server` 包装了 `Client`：

```go
type Service struct {
	*Client

	// handle Func
	HandleRequest func(msg *Call, wf ResultFunc)
}
```

`HandleRequest` 函数是处理调用的函数，在 `Service` 创建时由 `App` 传入，当 `Server` 收到一个请求时，进行处理后便调用该函数交由 `App` 层处理。

#### ResultFunc

`App` 层收到一个请求后，会创建一个 `ResultFunc`，里面存放着请求的 `请求ID` 和请求的 `连接ID`。并且 `App` 层在处理完数据后会调用 `ResultFunc.Result()` 向连接写回数据。

```go
// ResultFunc 回传结果的结构体
// 用于app层接收消息后，需要返回结果时调用其中的Result函数
type ResultFunc struct {

	// Result 返回结果的函数
	Result func(result Calls)

	// ConnID 请求的连接id
	ConnID string

	// ReqID 请求id
	ReqID string
}
```

#### 事件注册

`Client` 端只需要等待响应，进行回调，但是 `Server` 端需要接收请求。`Server` 层有 `Client` 层注入的处理函数，同样 `Logic` 也会向 `Dispatch` 层注入处理函数。

`Logic` 提供了 `Hook()` 函数和 `Handle()` 函数，用于扩展 `Logic` 的功能。

```go
func (c *Client) Hook(typ string, hookFunc interface{})

func (c *Client) Handle(typ string, handleFunc interface{})
```

目前这两个函数还未提供任何支持，仅仅是提供 `Dispatch` 层的 `Hook` 和 `Handle` 的转发。

## Dispatch层

`Dispatch` 负责通讯，通过为一个逻辑连接(例如一个连接池，中间的连接都是对一个地址发起的，会抽象为一个逻辑连接。)赋予一个唯一的 `连接ID`，隔离不同网络模式下的区别。上层只需要知道向哪一个逻辑连接发送数据包即可。

```go
// Dispatcher 通讯层的对外暴露的接口
type Dispatcher interface {

	// Link 连接到某个服务或中心
	// 会直接连接到指定的地址，[error]是用来返回连接时候的错误值的。
	// 连接断开不会在这里返回错误，而是提供一个hook，通过hook "close" 来捕获断开连接
	Link(addr string) error

	// ReLink 重新连接
	// 需要先调用 Link 之后才能调用ReLink，相当于是重新调用了一次Link，返回这次重连是否成功
	ReLink() bool

	// Send 发送一个帧
	// 发送一个帧出去，在不同的集群模式下有不同的表现
	// - default:
	// 发送到服务中心
	// - other:
	// 未实现
	Send(frame frame.Frame) error

	// SendTo 发送帧到指定连接
	SendTo(connID string, f frame.Frame) error

	// Listen 对一个地址开始监听
	Listen(addr string)

	// Close 释放资源
	Close()

	// Hook 对某些地方进行hook
	// 目前可以hook的：
	// - close
	Hook(typ string, hookFunc interface{})

	// Handle 对某些地方添加一个handle func来处理一些情况。
	// example:
	// dp.Handle("request",func(f *frame.Response) { fmt.Println(f) })
	// 目前可以handle的：
	// - client.handleResponse (response)
	// - client.handleRequest  (request)
	Handle(typ string, handleFunc interface{})

}
```

### Frame

Begonia 在网络中传输的数据帧格式是自定义的。分为 `header` 和 `payload`，除了 `payload` 之外的部分都是 `header` ：

```go
	    4      4       8    0 || 16  [  Len  ]
	{opcode}{version}{Len}{extendLen}{payload}
```

#### header

`header` 中第一个 `byte` 分为两个 4bit 部分。

**opcode**

`opcode` 的第一个 bit 是 `typcode`，代表这是一个请求帧还是一个响应帧。然后跟者三位的 `ctrlCode` 控制码，一共有 7 个，默认为 0b000，保留 6 个用于集群中的控制信号(比如心跳包等)。

**version**

版本号，代表当前的服务版本，默认为 0b0000，一共有 15 个版本。用于在某个服务更新后依旧需要提供旧版的服务所用。

通过简单位运算可以从 `byte` 上得到三个不同的 `code`：

```go
	typCode := 1 // 0 ~ 1
	dispatchCode := 4 // 0 ~ 7
	version := 8 // 0 ~ 15

	opcode := ((typCode<<3)|dispatchCode)<<4 | version

	version = opcode & 0b00001111
	dispatchCode = opcode >> 4 & 0b0111	
	typCode = opcode >> 7
	
	fmt.Printf("opcode: %08b %d\n", opcode, opcode)
	fmt.Printf("typCode:%01b %d\n", typCode, typCode)
	fmt.Printf("dispatchCode:%03b %d\n", dispatchCode, dispatchCode)
	fmt.Printf("versionCode:%04b %d\n", version, version)
----------------------------------------
opcode: 11001000 200
typCode:1 1
dispatchCode:100 4
versionCode:1000 8
```

**Len & ExtendLen**

第二个 `byte` 是 `len`，代表 `payload` 的长度。

这一位最大是 255，当它小于 255 时，它的数字就是后面跟着的剩余数据 `payload` 的长度。

当数据长度大于 255 时，这一位为 255 ，然后继续读取接下来两个 `byte`，这两个 `byte` 是 `extendLen`，当包长度小于 255 时不存在。

这两位用来表示最长 255 * 255 长度的包，这是 begonia 支持的最大的包长度。

#### payload

负载的设计比较简单：

```go
请求帧：
{reqId}0x00{service}0x00{fun}0x00{param}
响应帧：
{reqId}0x00{error}0x00{data}
```

由于reqID、service、fun、error都是字符串，并且里面的值都是可视值，所以直接使用了0x00作为分隔符。

## End

感谢阅读，这是一份写的非常仓促的设计文档。非常欢迎提出建议或批评。

如感兴趣也非常欢迎加入到开发中。目前正在筹划其他语言sdk的开发。

建议和批评请移至issue或添加我的WeChat：kkkkieran

