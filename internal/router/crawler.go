package router

// setupCrawlerRoutes 爬虫相关路由
func (r *Router) setupCrawlerRoutes() {
	api := r.engine.Group("/api")
	crawler := api.Group("/crawler")
	{
		crawler.POST("/start", r.handlers.Crawler.HandleCrawl)
		// 可以添加更多爬虫相关路由
		// crawler.GET("/status", r.handlers.Crawler.HandleStatus)
		// crawler.POST("/stop", r.handlers.Crawler.HandleStop)
	}
}
