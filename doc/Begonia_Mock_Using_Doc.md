# Begonia Mock 使用文档

## 初始化

如要获取一个带有mock功能的客户端，只需要调用begonia.NewClientWithMock函数来获取，其入参与获取默认客户端的参数一致。
代码示例如下：

```go
client := begonia.NewClientWithMock()
```

## 注册mock

注册mock函数与注册远程服务的方式是完全一样的，规则也是一样的：

```gomC
type testReceiver struct{}

func (r *testReceiver) NotVariadicISUFFunc(i int, s string, u uint, f float64) (x int) {
	return i + 1
}

func (r *testReceiver) VariadicFunc(s string, ints ...int) (x bool) {
	return len(ints) > 0
}

// 注册testReceiver接口体下的NotVariadicISUFFunc和VariadicFunc方法为mock函数
client.RegisterMock("mock", &testReceiver{}, "NotVariadicISUFFunc", "VariadicFunc")
```

注册mock函数成功后，调用该mock函数都将转换成对该方法的调用。

## 调用函数

获取到的客户端使用方式与默认客户端的使用方式一致。通过该函数获取的客户端将优先使用mock函数，若找到相应的mock函数会使用远程函数。

如果你想mock客户端在未找到mock函数时不会去寻找远程函数，或者你不想该mock客户端基于一个可用的客默认户端，可以调用begonia.NewMockClient函数来获取。

```go
client := begonia.NewMockClient()
```

## 自定义Mocker

mock客户端根据 `服务(service)` 分组操作mock选项，要操作一个service的mock，首先需要获取对应service的mock对象：

```go
// 获取服务名为 loop 的mock对象
mocker := client.GetServiceMocker("loop")
```

之后我们就可以对这个服务进行操作：

```go
// 注册mockLoop结构体下所有公开函数为mock函数原型
mocker.Register(&mockLoop{})

// 为 GetUid 函数增加自定义规则
mocker.Except("GetUid", []interface{}{"aaa"}, []interface{}{"2333"})
```

所有的mock函数都必须先注册，然后再增加自定义规则。注册只是获取mock函数的**函数类型**，用于后续自定义规则的合法判定。

自定义规则的函数签名如下：

```go
Except(funcName string, params []interface{}, out []interface{}) error
```

第一个参数是要增加规则的mock函数名字，第二个参数是入参期望（采用结构体注册的mock函数不需要receiver），第三个是出参期望。规则的优先级与增加先后顺序有关，越先注册的优先级越高。

> 这里的使用与gomock的Except使用相似

入参期望为一个空接口切片，每个函数的期望入参从左到右，与mock函数入参顺序一致。每个入参为定值，或为Matcher接口：

```go
// 为 GetUid 函数增加自定义规则，如果入参为("aaa", 23)，则为满足该规则，并出参("2333", true)
mocker.Except("GetUid", []interface{}{"aaa", 23}, []interface{}{"2333", true})

// 如果第一个参数为"aaa"，则为满足该规则
// 这里的第二个参数为 mock.NewAnyMatch() ，是begonia/app/mock包里的一个函数调用，其返回一个Match接口，表示任何参数都满足
// mock包里有多种返回Matcher接口的函数供使用，其提供不同的匹配规则
mocker.Except("GetUid", []interface{}{"aaa", mock.NewAnyMatch()}, []interface{}{"2333", true})
```

如果需要多个参数综合判断，则只需要这样：

```go
// 不管注册mock函数签名入参个数为多少，入参期望只有一个
mocker.Except("GetUid", []interface{}{
    mock.NewCustomMatch(func(s string, i int) bool {
        // logic
    }),
}, []interface{}{"2333", true})
```

mock.NewAnyMatch函数入参是一个空接口，返回一个Match接口。该函数的入参是一个**入参**与mock函数签名一致（采用结构体注册的mock函数不需要receiver），出参只是bool的函数。

如果你希望根据入参自定义出参，只需要在出参期望里**只传入一个mock.RetFunc**，如果匹配了该规则其会调用该函数并将返回值作为出参，例子如下：

```go
// 不管注册mock函数签名出参个数为多少，出参期望只有一个
// 当应用该规则时，会将入参传入该函数，然后你就可以根据入参自定义出参了
mocker.Except("GetUid", []interface{}{"aaa", 123}, 
    []interface{}{mock.RetFunc(func(params ...interface{}) (rets []interface{}, err error) {
        // logic
})})
```


