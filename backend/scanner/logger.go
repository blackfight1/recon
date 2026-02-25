package scanner

import (
	"fmt"
	"log"
	"recon-platform/database"
	"recon-platform/models"
)

// Logger 扫描日志记录器
type Logger struct {
	TaskID uint
}

// NewLogger 创建新的日志记录器
func NewLogger(taskID uint) *Logger {
	return &Logger{TaskID: taskID}
}

// Info 记录信息日志
func (l *Logger) Info(step, message string) {
	l.log("info", step, message)
	log.Printf("[Task %d] [%s] %s", l.TaskID, step, message)
}

// Success 记录成功日志
func (l *Logger) Success(step, message string) {
	l.log("success", step, message)
	log.Printf("[Task %d] [%s] ✓ %s", l.TaskID, step, message)
}

// Warning 记录警告日志
func (l *Logger) Warning(step, message string) {
	l.log("warning", step, message)
	log.Printf("[Task %d] [%s] ⚠ %s", l.TaskID, step, message)
}

// Error 记录错误日志
func (l *Logger) Error(step, message string) {
	l.log("error", step, message)
	log.Printf("[Task %d] [%s] ✗ %s", l.TaskID, step, message)
}

// log 内部日志记录方法
func (l *Logger) log(level, step, message string) {
	scanLog := &models.ScanLog{
		TaskID:  l.TaskID,
		Level:   level,
		Step:    step,
		Message: message,
	}
	database.DB.Create(scanLog)
}

// UpdateProgress 更新任务进度
func (l *Logger) UpdateProgress(progress int, currentStep string) {
	database.DB.Model(&models.ScanTask{}).
		Where("id = ?", l.TaskID).
		Updates(map[string]interface{}{
			"progress":     progress,
			"current_step": currentStep,
		})
}

// UpdateStatus 更新任务状态
func (l *Logger) UpdateStatus(status string) {
	database.DB.Model(&models.ScanTask{}).
		Where("id = ?", l.TaskID).
		Update("status", status)
}

// Infof 格式化信息日志
func (l *Logger) Infof(step, format string, args ...interface{}) {
	l.Info(step, fmt.Sprintf(format, args...))
}

// Successf 格式化成功日志
func (l *Logger) Successf(step, format string, args ...interface{}) {
	l.Success(step, fmt.Sprintf(format, args...))
}

// Warningf 格式化警告日志
func (l *Logger) Warningf(step, format string, args ...interface{}) {
	l.Warning(step, fmt.Sprintf(format, args...))
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(step, format string, args ...interface{}) {
	l.Error(step, fmt.Sprintf(format, args...))
}
