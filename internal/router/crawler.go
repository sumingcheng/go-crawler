package router

// setupCrawlerRoutes 爬虫相关路由
func (r *Router) setupCrawlerRoutes() {
	api := r.engine.Group("/api")
	crawler := api.Group("/crawler")
	{
		crawler.POST("/zhihu", r.handlers.Crawler.HandleCrawl)
	}
}
