package scheduler

import (
	"fmt"
	"log"
	"recon-platform/config"
	"recon-platform/database"
	"recon-platform/models"
	"recon-platform/scanner"

	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

// StartScheduler 启动定时任务调度器
func StartScheduler() {
	cronScheduler = cron.New()

	// 根据配置的间隔时间设置定时任务
	interval := config.AppConfig.Scanner.Interval
	cronExpr := getCronExpression(interval)

	log.Printf("Setting up scan scheduler with interval: %d hours (cron: %s)", interval, cronExpr)

	// 添加定时任务
	cronScheduler.AddFunc(cronExpr, func() {
		log.Println("Starting scheduled scan...")
		scanAllTargets()
	})

	// 启动调度器
	cronScheduler.Start()
	log.Println("Scheduler started successfully")

	// 启动时立即执行一次扫描（可选）
	go func() {
		log.Println("Running initial scan...")
		scanAllTargets()
	}()
}

// getCronExpression 根据小时数生成 cron 表达式
func getCronExpression(hours int) string {
	// 每 N 小时执行一次
	return fmt.Sprintf("0 */%d * * *", hours)
}

// scanAllTargets 扫描所有启用的目标
func scanAllTargets() {
	var targets []models.Target
	database.DB.Where("enabled = ?", true).Find(&targets)

	if len(targets) == 0 {
		log.Println("No enabled targets found")
		return
	}

	log.Printf("Found %d enabled targets to scan", len(targets))

	for _, target := range targets {
		go func(t models.Target) {
			if err := scanner.ScanTarget(t.ID, t.Domain); err != nil {
				log.Printf("Failed to scan target %s: %v", t.Domain, err)
			}
		}(target)
	}
}

// StopScheduler 停止调度器
func StopScheduler() {
	if cronScheduler != nil {
		cronScheduler.Stop()
		log.Println("Scheduler stopped")
	}
}
