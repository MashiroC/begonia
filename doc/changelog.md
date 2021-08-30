# Begonia change log (zh-CN)

## v1.1.0, 30 Aug 2021

- 修复了一个由于客户端发送请求后下线，导致服务中心代理器找不到返回的客户端的connID的`unknow bu ok error`
- 现在服务会在初始化时尝试拉取一次服务，如果不成功，会启动一个协程去持续拉，直到拉取成功，在这之间客户端调用函数会持续返回error
- 添加了初始化时连接失败的重试逻辑
- 添加了一个新的工具包`retry`，封装了重试的逻辑。

## v1.0.4, 28 Aug 2021

- 修改了注册的逻辑，现在会在注册时会hook link来在重连时发送注册信息

## v1.0.3, 27 Aug 2021

* 重新设计了`dispatch`的接口，现在`dispatch`唯一的接口是`Start()`
* 更改了dispatch层hook中`start`为`link`。该hook会在连接建立后启动
* 添加了grpc的demo用来做性能的比较
* 进行了大量的性能优化，包括：
  * 重构了frame的反序列化，现在不会大量使用`bufio.ReadBytes`去反复创建新的byte切片了。
  * 重构了`ResultFunc`，现在会直接在接收到数据后创建`context`，没有了臃肿的`ResultFunc`
  * 用`net.Conn`直接写替代了`bufio`写
* 正在尝试换一个高性能的网络库，目前所有性能损耗都在网络库上了。

## v1.0.2, 23 Aug 2021

- 因为现在全局用一个mode去判断begonia启动的模式，会导致在一个进程中启动多个不同mode的begonia节点时出现bug，所以添加了一个新的配置`option.Mode()`来指定启动的Mode
- 使用了一个`emptyCoder`替换了之前解析入参和出参为空时依旧会生成一堆schema
- 添加了reflect模式和ast模式的单元测试

## v1.0.1, 3 Aug 2021

- 给server这边的函数执行时的recover加上了日志记录
- 修复了一个远程函数拥有error时的解析错误的问题
- 现在ast模式也可以recover了


## v1.1.0, 25 Aug 2021

- 开始写change log了
- 修复了了自动生成工具中tag生成错误的bug