package proxy

// 即使第三方的包，通过golang 的 duck type，我们也能抽象出一套接口。
// 不过如果是接口众多，我们可以考虑用 proxy 直接继承、组合使用第三方包
