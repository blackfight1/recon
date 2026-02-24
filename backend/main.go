package main

import (
	"fmt"
	"log"
	"recon-platform/config"
	"recon-platform/database"
	"recon-platform/router"
	"recon-platform/scheduler"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 启动定时任务调度器
	scheduler.StartScheduler()

	// 启动 HTTP 服务器
	r := router.SetupRouter()
	port := config.AppConfig.Server.Port

	log.Printf("Server starting on port %d...", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
