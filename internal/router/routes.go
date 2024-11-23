package router

func (r *Router) SetupRoutes() {
	// 注册业务路由
	r.setupCrawlerRoutes()
	// 注册系统路由（健康检查等）
	r.setupHealthRoutes()
}
