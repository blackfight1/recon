package router

import (
	"recon-platform/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API 路由组
	api := r.Group("/api")
	{
		// 目标管理
		targets := api.Group("/targets")
		{
			targets.GET("", controllers.GetTargets)
			targets.POST("", controllers.CreateTarget)
			targets.GET("/:id", controllers.GetTarget)
			targets.PUT("/:id", controllers.UpdateTarget)
			targets.DELETE("/:id", controllers.DeleteTarget)
			targets.POST("/:id/scan", controllers.TriggerScan)
		}

		// 子域名
		subdomains := api.Group("/subdomains")
		{
			subdomains.GET("", controllers.GetSubdomains)
			subdomains.GET("/target/:target_id", controllers.GetSubdomainsByTarget)
		}

		// 变更日志
		changes := api.Group("/changes")
		{
			changes.GET("", controllers.GetChanges)
			changes.GET("/target/:target_id", controllers.GetChangesByTarget)
		}

		// 扫描任务
		tasks := api.Group("/tasks")
		{
			tasks.GET("", controllers.GetTasks)
			tasks.GET("/:id", controllers.GetTask)
			tasks.GET("/:id/logs", controllers.GetTaskLogs)
			tasks.GET("/:id/progress", controllers.GetTaskProgress)
		}

		// 快速扫描
		api.POST("/quick-scan", controllers.QuickScan)

		// 统计信息
		api.GET("/stats", controllers.GetStats)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
