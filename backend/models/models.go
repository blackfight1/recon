package models

import (
	"time"
)

// Target 监控目标（主域名）
type Target struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Domain      string    `gorm:"uniqueIndex;not null" json:"domain"`
	Description string    `json:"description"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Subdomain 子域名记录
type Subdomain struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TargetID   uint      `gorm:"index;not null" json:"target_id"`
	Subdomain  string    `gorm:"index;not null" json:"subdomain"`
	Status     string    `gorm:"default:'unknown'" json:"status"` // alive, dead, new
	StatusCode int       `json:"status_code"`
	Title      string    `json:"title"`
	Server     string    `json:"server"`
	IP         string    `json:"ip"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ScanTask 扫描任务
type ScanTask struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	TargetID    uint       `gorm:"index;not null" json:"target_id"`
	TaskType    string     `json:"task_type"`                       // subdomain, port, screenshot, vuln
	Status      string     `gorm:"default:'pending'" json:"status"` // pending, running, completed, failed
	Result      string     `gorm:"type:text" json:"result"`
	ErrorMsg    string     `gorm:"type:text" json:"error_msg"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ChangeLog 变更日志
type ChangeLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TargetID   uint      `gorm:"index;not null" json:"target_id"`
	ChangeType string    `json:"change_type"` // subdomain_new, subdomain_dead, subdomain_alive
	Content    string    `gorm:"type:text" json:"content"`
	Notified   bool      `gorm:"default:false" json:"notified"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 指定表名
func (Target) TableName() string {
	return "targets"
}

func (Subdomain) TableName() string {
	return "subdomains"
}

func (ScanTask) TableName() string {
	return "scan_tasks"
}

func (ChangeLog) TableName() string {
	return "change_logs"
}
