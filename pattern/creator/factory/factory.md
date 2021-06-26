---
title: pattern-工厂模式
tags: pattern
notebook: pattern
---

工厂模式分为：`工厂方法模式`和`抽象工厂模式`，我们常说的`简单工厂模式`也被认为是`工厂方法模式`的特例。
`简单工厂模式`最常用到，`工厂方法模式`偶尔，`抽象工厂模式`我几乎没有见过
由于，多数人都习惯把工厂模式区分为`简单工厂模式`、`工厂方法模式`、`抽象工厂模式`，所以接下来我如此去介绍工厂模式相关的内容了。

## 简单工厂
简单工厂模式是经常用到的，比如：
```Go
def get_duck(color):
    if color == "yello":
        return YelloDuck()
    elif color == "blue":
        return BlueDuck()
    else:
        return DefaultDuck()
```
这就是传说中的简单工厂模式，不过，在Go里我们一般不会返回多个struct，而是返回一个interface，而具体实现，都满足这个interface， 比如，如果我们做一个短信服务，肯定要对接多个短信渠道，比如阿里云、腾讯云，那么就可以这样：

```go
type SMSSender interface {
	Send(content string, receivers []string) error
}

type AliyunSMS struct{}

func (a *AliyunSMS) Send(content string, receivers []string) error {
	// pass
}

type TencentSMS struct{}

func (t *TencentSMS) Send(content string, receivers []string) error {
	// pass
}

// 简单工厂在这里
func GetSMSSender(channel string) SMSSender {
	if channel == "aliyun" {
		return &AliyunSMS{}
	} else if channel == "tencent" {
		return &TencentSMS{}
	} else {
		// 略
	}
}
```
但是这只是一种用法，还有一种写法上的变种：
```go
type SMSSender interface {
	Send(content string, receivers []string) error
}

type AliyunSMS struct{}

func (a *AliyunSMS) Send(content string, receivers []string) error {
	return nil
}

type TencentSMS struct{}

func (t *TencentSMS) Send(content string, receivers []string) error {
	return nil
}

var senderMapper = map[string]SMSSender{
	"aliyun":  &AliyunSMS{},
	"tencent": &TencentSMS{},
}

// 简单工厂在这里
func GetSMSSender(channel string) SMSSender {
	sender, exist := senderMapper[channel]
	if !exist {
		// 略
	}

	return sender
}
```
使用

```golang
func SendEmail(content string, receivers []string) error {
    channel := parseChannel(content)
    sender := SMSSenderFactory.getSMSSender(channel)
    sender.Send(content, receivers)
}
```

当然，这里的区别在于，使用一个mapper之后，就节省了一堆的 if...else...，不过缺点就是并非每次都 实例化了对应的sender，当然也是可以通过反射做到的，不过不推荐，所以实际上用哪种 写法，还是要结合实际情况来看。

## 工厂方法
假设我们上面的发送短信服务依赖了一下对象，初始化很麻烦，继续使用简单工厂模式会是怎样呢。
```golang
// 工厂方法
func GetSMSSender(channel string) SMSSender {
	if strings.toLower(channel) == "aliyun" {
        aliyunUser := aliyunSdk.NewUser(globalConfig.UserKey)
        opts := &aliyunSdk.ConnOptions{
            // ...
        }
        return NewAliyunSMS(aliyunUser, opts)
    }else if {
        // ...
    }

	return sender
}
```
当每个SMS对象的组装逻辑都很复杂，那整个`GetSMSSender`工厂方法就会显得臃肿，复杂，不利于后续的拓展、维护。

## 是简单工厂还是工厂方法

之所以将某个代码块剥离出来，独立为函数或者类，原因是这个代码块的逻辑过于复杂，剥离之后能让代码更加清晰，更加可读、可维护。但是，如果代码块本身并不复杂，就几行代码而已，我们完全没必要将它拆分成单独的函数或者类。
基于这个设计思想，当对象的创建逻辑比较复杂，不只是简单的 new 一下就可以，而是要组合其他类对象，做各种初始化操作的时候，我们推荐使用工厂方法模式，将复杂的创建逻辑拆分到多个工厂类中，让每个工厂类都不至于过于复杂。
而如果只是很简单的创建对象，那么直接使用简单工厂模式。
## 抽象工厂
工厂模式还有一种，是抽象工厂模式，这个似乎不太常用，至少我没有在代码很少遇到。

### 实战
> 接下来我们来看一个复杂的场景，每个工厂生产对象初始化入参不同怎么办

现实场景
现在我们要适配一个名为OPS的服务，该服务主要功能是向OPS发送心跳、从OPS拉取配置以及通过OPS下发更新包。假设目前有2个场景，一个基础场景，一个特殊场景。

接口定义
首先，我们将公共的接口抽出来，定义为一个公共的接口。例如目前的场景下，定义接口OpsImpl如下

type OpsImpl interface {
    SendHeartbeat() error
    DoUpdate() error
    DoConfigUpload() error
}
接口实现
接下来，我们针对基础场景和特殊场景，创建struct实现接口。首先是基础场景baseOps实现如下：

type baseOps struct {
    postUrl string
}

func (base *baseOps) SendHeartbeat() error {
    fmt.Println("BaseOps: Send heartbeat")
    fmt.Println("BaseOps: send to url: ", base.postUrl)
    return nil
}
func (base *baseOps) DoUpdate() error {
    fmt.Println("BaseOps: DoUpdate")
    return nil
}
func (base *baseOps) DoConfigUpload() error {
    fmt.Println("BaseOps: DoConfigUpload")
    return nil
}
接下来是特殊场景，特殊场景是基础场景的子类，我们通过组合的模式实现该逻辑，其中特殊场景的更新逻辑与基础场景保持一致，因此我们可以省略其实现，复用父类的实现，而心跳和配置拉取逻辑不一致，需要重新写业务逻辑，具体实现如下：

type specialOps struct {
    baseOps
    sendConfig bool
}

func (ops *specialOps) SendHeartbeat() error {
    fmt.Println("SpecialOps: SendHeartbeat")
    return nil
}
func (ops *specialOps) DoConfigUpload() error {
    fmt.Println("SpecialOps: DoConfigUpload")
    if ops.sendConfig {
        fmt.Println("SpecialOps: upload config")
    } else {
        fmt.Println("SpecialOps: no need to upload config")
    }
    return nil
}

// not implement DoUpdate, same to baseOps
工厂函数
接下来，我们为不同的实现工厂方法，返回同样的接口。这些工厂方法接受同样的参数，我们将参数通过map[string]interface{}进行传递。实现如下：

type OpsFactory func(conf map[string]interface{}) (OpsImpl, error)

func NewBaseOps(conf map[string]interface{}) (OpsImpl, error) {
    fmt.Println("BaseOps: Create")
    postUrl, ok := conf["PostUrl"]
    if !ok {
        return nil, errors.New("[postUrl] has not been set in config map")
    }
    return &baseOps{
        postUrl: postUrl.(string),
    }, nil
}

func NewSpecialOps(conf map[string]interface{}) (OpsImpl, error) {
    fmt.Println("specialOps: Create")
    sendConfig, ok := conf["SendConfig"]
    if !ok {
        return nil, errors.New("[SendConfig] has not been set in config map")
    }
    return &specialOps{
        sendConfig: sendConfig.(bool),
    }, nil
}
注册工厂
接下来，我们通过一个公共的函数RegisterOps来注册需要用到的工厂，并且通过初始化函数init在程序启动前注册这两个工厂。具体实现如下：

var opsFactories = make(map[opsTypeEnum]OpsFactory)

func RegisterOps(opsType opsTypeEnum, factory OpsFactory) {
    if factory == nil {
        panic(fmt.Sprintf("Ops factory %s does not exist", string(opsType)))
    }
    _, ok := opsFactories[opsType]
    if ok {
        fmt.Printf("Ops factory %s has been registered already\n", string(opsType))
    } else {
        fmt.Printf("Register ops factory %s\n", string(opsType))
        opsFactories[opsType] = factory
    }
}

func init() {
    RegisterOps(BaseType, NewBaseOps)
    RegisterOps(SpecialType, NewSpecialOps)
}
创建工厂
最后我们通过以下函数，就可以方便的创建工厂，返回对应的Ops接口。

func CreateOps(conf map[string]interface{}) (OpsImpl, error) {
    opsType, ok := conf["OpsType"]
    if !ok {
        fmt.Println("No ops type in config map. Use base ops type as default.")
        opsType = BaseType
    }
    OpsFactory, ok := opsFactories[opsType.(opsTypeEnum)]
    if !ok {
        availableOps := make([]string, len(opsFactories))
        for k, _ := range opsFactories {
            availableOps = append(availableOps, string(k))
        }
        return nil, errors.New(fmt.Sprintf("Invalid ops type. Must be one of: %s", strings.Join(availableOps, ", ")))
    }
    fmt.Println("Create ops: ", opsType)
    return OpsFactory(conf)
}
测试函数
最终，我们可以通过以下方式简单的创建不同场景的接口。

func main() {
    baseOps, err := ops.CreateOps(map[string]interface{}{
        "OpsType": ops.BaseType,
        "PostUrl": "http://ops.cloud.com/send_heartbeat",
    })
    if err != nil {
        fmt.Println("create baseOps failed, err: ", err.Error())
        return
    }
    baseOps.DoConfigUpload() // Output: BaseOps: DoConfigUpload
    baseOps.DoUpdate()       // Output: BaseOps: DoUpdate
    baseOps.SendHeartbeat()  // Output: BaseOps: Send heartbeat

    specialOps, err := ops.CreateOps(map[string]interface{}{
        "OpsType":    ops.SpecialType,
        "SendConfig": true,
    })
    if err != nil {
        fmt.Println("create specialOps failed, err: ", err.Error())
        return
    }
    specialOps.DoConfigUpload() // Output: SpecialOps: DoConfigUpload
    specialOps.DoUpdate()       // Output: BaseOps: DoUpdate
    specialOps.SendHeartbeat()  // Output: SpecialOps: SendHeartbeat
}