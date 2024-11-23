package router

// SetupRoutes 注册所有路由
func (r *Router) SetupRoutes() {
	// 注册系统路由（健康检查等）
	r.setupHealthRoutes()
	// 注册业务路由
	r.setupCrawlerRoutes()
}
