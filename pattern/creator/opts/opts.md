自己的项目调用这两个包时，有如下的代码（如果有兴趣可以自己去github上看完整的源码）
1、用elastic包构建聚合查询

    agg := elastic.NewTermsAggregation().
        Size(10000).
        Field("data.sce")
2、用grpc构造连接

conn, err := grpc.DialContext(connCtx, "127.0.0.1:5555", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        return err
    }
抛开两个函数的功能不谈，其实两个函数都可以理解为一种构造或者初始化函数，其Size，Field，WithInsecure，WithBlock等函数都是在初始化或者构造时提供某种参数而已。

那我们提炼一下，就是当我们构造一个对象时，一般会提供很多参数用来构造，但是不同场景，或者不同条件下，需要的参数又不同，如何来封装这些构造函数来方便使用呢？

以如下结构为例，假设是一个代理

type Agency struct {
 IP       int32         //required
 Protocol string        //optional
 Timeout  time.Duration //optional
}
方法一：
直接构造不同的构造函数：

func NewAgency(ip int32) *Agency
func NewAgencyWithProtocol(ip int32, p string) *Agency
func NewAgencyWithProtocolAndTimeout(ip int32, p string, t time.Duration) *Agency
......
你会发现需要定义一系列的函数，而且随着参数增多，可能扩充的函数也比较多，同时由于go中不支持多态，每个函数还要有不同的名称，当你这次调用了A，下次增加参数时，还得改为调用B，可见麻烦多多，既不好看，也不好用，还不好维护。当然实际项目中可能也没人会这么写，这里只是举例而已。

方法二：
直接将Agency结构体作为参数，这样参数不就固定了吗，而且一举解决所有问题，只用一个构造函数即可。

func NewAgency(a Agency) *Agency
这种方法现实项目中确实也有使用的哦，那对于这种简单的结构体，还算比较方便简洁，但是如果结构体成员较多（像goroutine等结构），动辄20+以上，那你初始化的时候还得先一一确定参数，然后再去调用构造函数是不是也麻烦，而且很多参数其实用不到赋值，只要默认值就够用了。而且成员变量一多，你可能都不知道那些是必选参数，那些是可选参数了。
那如何解决这个问题呢？

方法三：
修改Agency的结构体，将可选成员和必选成员分开：

type AgencyOption struct {
    Protocol string        //optional
    Timeout  time.Duration //optional
}
type Agency struct {
    IP       int32         //required
    AgencyOption
}
然后定义一个构造函数

func NewAgency(ip int32, param *AgencyOption) *Agency
这个函数区分了必选项和可选项，ip必填，param可选，如果param为nil则不用对可选参数赋值。
这种方案相对于上一个方法略有改进，对于必选参数一目了然，但是对于参数较多的场景还是没有根本解决。

方法四：
直接将可选参数不放在构造函数中，定义多个设置函数，例如：

func (a *Agency)SetProtocol(p string) {
   a.Protocol = p
}
func (a *Agency)SetTimeout(t time.Duration) {
   a.Timeout = t
}
调用者使用方式：

a1 := NewAgency(ip)
a1.SetProtocol("udp")
a1.SetTimeout(100)
方法五：
上述方法需要多次调用，我们做个修改：

func (a *Agency)SetProtocol(p string) *Agency{
   a.Protocol = p
   return a
}
func (a *Agency)SetTimeout(t time.Duration) *Agency{
   a.Timeout = t
   return a
}
这样，调用者可以使用链式调用：

a1 := NewAgency(ip).SetTimeout(100).SetProtocol("udp")
这也是一种常用的方法。
像文章开篇提到的elastic库就是这么玩的，每次查询的时候可能有很多参数要设置，直接连续调用即可，很清晰，这也是从这个开源库中学到。以后就可以直接应用到实际项目中了哦！
那如果就想将参数一把传入，一次性初始化完成呢？

方法六：
一次性传入任意多个参数，首先我们可以想到go支持变参，例如func Add(base int, others ...int)，可以处理任意个数的int类型，但是我们的参数一般是不一样的，那我们如何利用这种方案呢？
我们大胆想象一下如果有这么一种通用类型可以使用(先不考虑返回值)：

func NewAgency(ip int32, options ...Option)
如果能实现这样一个构造函数，那么可选参数的问题也就搞定了！！
但是这个通用的Option如何定义呢？？使用某一种具体类型肯定是不能完成的，那是否可以将这个Option定义为一个函数或者接口类型呢？

先从函数思考，看能否实现：
一个函数，无非是函数名，入参，逻辑处理，返回值这些东东
函数名，whatever，随便起个能自注释的即可；
入参，先放一下；
逻辑处理，就是这个函数要做啥，想想我们的最终目的就是将可选参数设置到我们的对象中去！！那么这个逻辑处理就类似于：

agency.Timeout = time
以及
agency.Protocol = "udp"
等等...
从逻辑处理看我们要将参数设置到对象中，那这个对象是不是可以作为我们的共同参数，那么这个Option的定义是不是可以为：

type Option func(a *Agency)
因为是要改变Agency中的值，所以用的是指针作为入参。

既然Option定义好了，那么针对每个参数我们来实现由参数如何转换为这种Option传入构造函数吧，即创建入参是参数，但是返回值是Option类型的函数
对于超时时间：

func Timeout(t time.Duration) Option{
   return func(a *Agency){
      a.Timeout = t
   }
}
对于协议设置：

func Protocol(p string) Option{
   return func(a *Agency){
      a.Protocol = p
   }
}
构造函数为：

func NewAgency(ip int32, options ...Option) *Agency {
   a := &Agency{IP:ip}
   for _, opt := range options{
      opt(a)
   }
   
   return a
}
调用者为：

a1 := NewAgency(ip)
a2 := NewAgency(ip, Timeout(100))
a4 := NewAgency(ip, Timeout(200), Protocol("udp"))
这样就清晰明了了吧，大功告成！

方法七：
上个方法中定义的Option是一个函数类型，那么接口类型是否也可以胜任呢？
答案也是可以的，这个是我从grpc的实现反向思考的过程，大家可以参考，或者直接撸grpc的源码看：)
我们还是以Agency为例
首先定义这个Option接口，由于在go中通常定义的接口名都带er，我们也遵照传统定义为：

type Optioner interface {
   apply()
}
接口先只定义了一个名字，入参和返回值待定。
既然有了接口，那么我们就要定义一个类型来实现这个接口：

type RealOption struct {
}
func (ro *RealOption)apply(){
}
雏形就有了，那么如何来填这些定义的内容呢？
先别急，我们继续定义设置参数的函数，返回值类型都要为Optioner，所以其模型类似为：

func SetProtocol(p string) Optioner {
   return &RealOption{
   }
}
func SetTimeout(t time.Duration) Optioner {
   return &RealOption{
   }
}
再继续看我们最终能够提供的构造函数，应该是这个样子：

func NewAgency(ip int32, options ...Optioner) *Agency {
   a := &Agency{IP:ip}
   for _, opt := range options{
      opt.apply()
   }
   return a
}
核心还是遍历可变参数options，去调用对应的接口设置相应的参数，从这里作为突破口，那么apply这个接口类型应该定义成什么呢，是不是呼之欲出了，只要增加个入参，无需返回值

apply(agency *Agency)
有了入参后，上述涉及参数的各个定义修改为：

func (ro *RealOption)apply(agency *Agency){
}
type Optioner interface {
   apply(agency *Agency)
}
func NewAgency(ip int32, options ...Optioner) *Agency {
   a := &Agency{IP:ip}
   for _, opt := range options{
      opt.apply(a)
   }
   return a
}
既然接口定义好了，那么看如何实现参数设置函数的逻辑，对于SetTimeout函数，目的是将入参t传到对象中去，那么就类似于：

func SetTimeout(t time.Duration) Optioner {
   return &RealOption{
      ?:t,
   }
}
如果?的位置是具体的类型，那么这个t其实是设置到了RealOption中，并没有设置到Agency中，那么怎么办呢，还记得上个方法中的Option定义吗 ，其就是一个通用类型的函数，将参数设置到Agency里，那么这个地方也用这种方式呢，即：

func SetTimeout(t time.Duration) Optioner {
   return &RealOption{
      ?:func(agency *Agency){
         agency.Timeout = t
      },
   }
}
那这样RealOption的定义也就出来了：

type RealOption struct {
   f func(agency *Agency)
}
那么所有的内容基本都完成了，整体代码如下：

type Optioner interface {
   apply(agency *Agency)
}
type RealOption struct {
   f func(agency *Agency)
}
func (ro *RealOption) apply(agency *Agency) {
   ro.f(agency)
}
func SetProtocol(p string) Optioner {
   return &RealOption{
      f: func(agency *Agency) {
         agency.Protocol = p
      },
   }
}
func SetTimeout(t time.Duration) Optioner {
   return &RealOption{
      f: func(agency *Agency) {
         agency.Timeout = t
      },
   }
}
func NewAgency(ip int32, options ...Optioner) *Agency {
   a := &Agency{IP: ip}
   for _, opt := range options {
      opt.apply(a)
   }
   return a
}
调用方式跟上一种方法一样，不过传入的可选参数是接口而已
继续优化一下，对RealOption结构也提供一个构造函数，那么参数设置函数改为：

func NewRealOption(f func(agency *Agency)) *RealOption {
   return &RealOption{
      f: f,
   }
}
func SetProtocol(p string) Optioner {
   return NewRealOption(func(agency *Agency) {
      agency.Protocol = p
   })
}
func SetTimeout(t time.Duration) Optioner {
   return NewRealOption(func(agency *Agency) {
      agency.Timeout = t
   })
}
OK,大功告成！这种方法就对应了开篇提到的grpc中实现的方法。