![20210623225318](https://raw.githubusercontent.com/litao-2071/static-res/master/docs-images/20210623225318.png)

更倾向非延迟加载，更早暴露问题；懒汉式虽然支持延迟加载，但是这只是把冷启动时间放到了第一次使用的时候，并没有本质上解决问题，并且为了实现懒汉式还不可避免的需要加锁