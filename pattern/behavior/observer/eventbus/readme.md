`异步非阻塞`除了能实现代码解耦之外，还能提高代码的执行效率；

进程间的观察者模式解耦更加彻底，一般是基于消息队列来实现，用来实现不同进程间的被观察者和观察者之间的交互。

这里我们聚焦于`异步非阻塞`的观察者模式，实现一个简单的`Event bus`。

功能更全的`Event bus`，可以查看：https://github.com/asaskevich/EventBus/blob/master/event_bus.go
