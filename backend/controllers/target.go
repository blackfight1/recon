package controllers

import (
	"net/http"
	"recon-platform/database"
	"recon-platform/models"
	"recon-platform/scanner"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTargets 获取所有目标
func GetTargets(c *gin.Context) {
	var targets []models.Target
	database.DB.Order("created_at desc").Find(&targets)
	c.JSON(http.StatusOK, gin.H{"data": targets})
}

// GetTarget 获取单个目标
func GetTarget(c *gin.Context) {
	id := c.Param("id")
	var target models.Target

	if err := database.DB.First(&target, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": target})
}

// CreateTarget 创建目标
func CreateTarget(c *gin.Context) {
	var input struct {
		Domain      string `json:"domain" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target := models.Target{
		Domain:      input.Domain,
		Description: input.Description,
		Enabled:     true,
	}

	if err := database.DB.Create(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create target"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": target})
}

// UpdateTarget 更新目标
func UpdateTarget(c *gin.Context) {
	id := c.Param("id")
	var target models.Target

	if err := database.DB.First(&target, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	var input struct {
		Domain      string `json:"domain"`
		Description string `json:"description"`
		Enabled     *bool  `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Domain != "" {
		target.Domain = input.Domain
	}
	if input.Description != "" {
		target.Description = input.Description
	}
	if input.Enabled != nil {
		target.Enabled = *input.Enabled
	}

	database.DB.Save(&target)
	c.JSON(http.StatusOK, gin.H{"data": target})
}

// DeleteTarget 删除目标
func DeleteTarget(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Target{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}

// TriggerScan 手动触发扫描
func TriggerScan(c *gin.Context) {
	id := c.Param("id")
	var target models.Target

	if err := database.DB.First(&target, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	// 异步执行扫描
	go scanner.ScanTarget(target.ID, target.Domain)

	c.JSON(http.StatusOK, gin.H{"message": "Scan triggered successfully"})
}

// GetSubdomains 获取所有子域名
func GetSubdomains(c *gin.Context) {
	var subdomains []models.Subdomain

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	status := c.Query("status")

	query := database.DB.Model(&models.Subdomain{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("last_seen desc").Limit(pageSize).Offset(offset).Find(&subdomains)

	c.JSON(http.StatusOK, gin.H{
		"data":  subdomains,
		"total": total,
		"page":  page,
	})
}

// GetSubdomainsByTarget 获取指定目标的子域名
func GetSubdomainsByTarget(c *gin.Context) {
	targetID := c.Param("target_id")
	var subdomains []models.Subdomain

	database.DB.Where("target_id = ?", targetID).
		Order("last_seen desc").
		Find(&subdomains)

	c.JSON(http.StatusOK, gin.H{"data": subdomains})
}

// GetChanges 获取所有变更日志
func GetChanges(c *gin.Context) {
	var changes []models.ChangeLog

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	var total int64
	database.DB.Model(&models.ChangeLog{}).Count(&total)

	offset := (page - 1) * pageSize
	database.DB.Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&changes)

	c.JSON(http.StatusOK, gin.H{
		"data":  changes,
		"total": total,
		"page":  page,
	})
}

// GetChangesByTarget 获取指定目标的变更日志
func GetChangesByTarget(c *gin.Context) {
	targetID := c.Param("target_id")
	var changes []models.ChangeLog

	database.DB.Where("target_id = ?", targetID).
		Order("created_at desc").
		Find(&changes)

	c.JSON(http.StatusOK, gin.H{"data": changes})
}

// GetTasks 获取所有任务
func GetTasks(c *gin.Context) {
	var tasks []models.ScanTask

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	database.DB.Model(&models.ScanTask{}).Count(&total)

	offset := (page - 1) * pageSize
	database.DB.Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&tasks)

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"total": total,
		"page":  page,
	})
}

// GetTask 获取单个任务
func GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.ScanTask

	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// GetStats 获取统计信息
func GetStats(c *gin.Context) {
	var stats struct {
		TotalTargets    int64 `json:"total_targets"`
		EnabledTargets  int64 `json:"enabled_targets"`
		TotalSubdomains int64 `json:"total_subdomains"`
		AliveSubdomains int64 `json:"alive_subdomains"`
		NewSubdomains   int64 `json:"new_subdomains"`
		RecentChanges   int64 `json:"recent_changes"`
	}

	database.DB.Model(&models.Target{}).Count(&stats.TotalTargets)
	database.DB.Model(&models.Target{}).Where("enabled = ?", true).Count(&stats.EnabledTargets)
	database.DB.Model(&models.Subdomain{}).Count(&stats.TotalSubdomains)
	database.DB.Model(&models.Subdomain{}).Where("status = ?", "alive").Count(&stats.AliveSubdomains)
	database.DB.Model(&models.Subdomain{}).Where("status = ?", "new").Count(&stats.NewSubdomains)
	database.DB.Model(&models.ChangeLog{}).Where("created_at > NOW() - INTERVAL '7 days'").Count(&stats.RecentChanges)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// QuickScan 快速扫描（不添加到监控列表）
func QuickScan(c *gin.Context) {
	var input struct {
		Domain string `json:"domain" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建临时任务（targetID 为 0）
	task := &models.ScanTask{
		TargetID:    0,
		Domain:      input.Domain,
		TaskType:    "subdomain",
		Status:      "pending",
		Progress:    0,
		CurrentStep: "等待开始",
	}
	if err := database.DB.Create(task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// 异步执行扫描（使用已创建的 task.ID）
	go func(taskID uint, domain string) {
		scanner.ScanWithTask(taskID, 0, domain)
	}(task.ID, input.Domain)

	c.JSON(http.StatusOK, gin.H{
		"message": "Quick scan started",
		"task_id": task.ID,
	})
}

// GetTaskLogs 获取任务日志
func GetTaskLogs(c *gin.Context) {
	taskID := c.Param("id")
	var logs []models.ScanLog

	database.DB.Where("task_id = ?", taskID).
		Order("created_at asc").
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetTaskProgress 获取任务进度
func GetTaskProgress(c *gin.Context) {
	taskID := c.Param("id")
	var task models.ScanTask

	if err := database.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}
