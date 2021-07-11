package middleware

import (
	"net/http"
	"net/http/httputil"
)

//让 ChainHandlerFunc 继承 http.Handler 方便做其他函数的参数
type ChainHandlerFunc func(rw http.ResponseWriter, req *http.Request)

func (f ChainHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

//中间件方法类型
type MiddleWareHandlerFunc func(next http.Handler) http.Handler

//让 WrapHandlerEntity 继承 http.Handler 方柏霓其他函数的参数
type WrapHandlerEntity struct {
	Handler http.Handler
}

func (w *WrapHandlerEntity) ServeHTTP(w2 http.ResponseWriter, r *http.Request) {
	w.Handler.ServeHTTP(w2, r)
}

//代理router, 用以创建链式结构
type ChainRouter struct {
	*httputil.ReverseProxy
	prev       *ChainRouter
	middleware MiddleWareHandlerFunc //中间支持
}

func NewChainRouter(proxy *httputil.ReverseProxy) *ChainRouter {
	return &ChainRouter{
		ReverseProxy: proxy,
	}
}

//创建 middleware 链式结构
func (p *ChainRouter) Use(middlewares ...MiddleWareHandlerFunc) *ChainRouter {
	if len(middlewares) == 0 {
		return p
	}

	router := p
	for _, mw := range middlewares {
		router = router.use(mw)
	}

	return router
}

//单步链式创建
//尾插法创建链表
func (p *ChainRouter) use(mw MiddleWareHandlerFunc) *ChainRouter {
	return &ChainRouter{
		prev:         p,
		ReverseProxy: p.ReverseProxy,
		middleware:   mw,
	}
}

//基于链表构建方法链
func (p *ChainRouter) genChainFunc(handler http.Handler) http.Handler {
	wraphandler := &WrapHandlerEntity{
		Handler: handler,
	}

	chain := handler
	router := p

	for router.prev != nil {
		if router.middleware != nil {
			//一次调用如下
			//通过调用倒数第一 middleware 的 MiddleWareHandlerFunc (初始化 http.Handler) 获取 http.Handler方法
			//通过调用倒数第二 middleware 的 MiddleWareHandlerFunc (上步所得 http.Handler) 获取 http.Handler 方法
			// ...
			//通过调用第一 middleware 的 MiddleWareHandlerFunc (上步所得 http.Handler) 获取 http.Handler
			//形成方法的嵌套, 先加入的在外边

			//形成的方法类似如下
			//ssh := func(h http.Handler, w http.ResponseWriter, r *http.Request) {
			//	//middleware 1 header
			//		//middleware 2 header
			//			//middleware 3 header
			//				h.ServeHTTP(w, r)
			//			//middleware 3 footer
			//		//middleware 2 footer
			//	//middleware 1 footer
			//}

			chain = router.middleware(wraphandler)
		}

		wraphandler = &WrapHandlerEntity{
			Handler: chain,
		}

		router = router.prev
	}

	return &WrapHandlerEntity{
		Handler: chain,
	}
}

//外部服务接口
func (p *ChainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// step 1 基于链表构建方法
	chainHandler := p.genChainFunc(p.ReverseProxy)
	// step 2 调用方法链

	chainHandler.ServeHTTP(w, r)
}
