package notifier

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"recon-platform/config"
	"recon-platform/models"
	"time"
)

// SendChangeNotification 发送变更通知
func SendChangeNotification(domain string, changes []models.ChangeLog) {
	if len(changes) == 0 {
		return
	}

	// 统计变更
	newCount := 0
	aliveCount := 0
	deadCount := 0

	for _, change := range changes {
		switch change.ChangeType {
		case "subdomain_new":
			newCount++
		case "subdomain_alive":
			aliveCount++
		case "subdomain_dead":
			deadCount++
		}
	}

	message := fmt.Sprintf("🔍 资产变更通知\n\n")
	message += fmt.Sprintf("目标: %s\n", domain)
	message += fmt.Sprintf("时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	if newCount > 0 {
		message += fmt.Sprintf("🆕 新增子域名: %d 个\n", newCount)
	}
	if aliveCount > 0 {
		message += fmt.Sprintf("✅ 恢复存活: %d 个\n", aliveCount)
	}
	if deadCount > 0 {
		message += fmt.Sprintf("❌ 失效子域名: %d 个\n", deadCount)
	}

	message += "\n详细信息:\n"
	for i, change := range changes {
		if i >= 10 {
			message += fmt.Sprintf("\n... 还有 %d 条变更，请登录系统查看", len(changes)-10)
			break
		}

		emoji := "🔹"
		switch change.ChangeType {
		case "subdomain_new":
			emoji = "🆕"
		case "subdomain_alive":
			emoji = "✅"
		case "subdomain_dead":
			emoji = "❌"
		}

		message += fmt.Sprintf("%s %s\n", emoji, change.Content)
	}

	// 发送到企业微信
	if config.AppConfig.Notification.Wecom.Enabled {
		sendWecom(message)
	}

	// 发送到钉钉
	if config.AppConfig.Notification.Dingtalk.Enabled {
		sendDingtalk(message)
	}
}

// sendWecom 发送企业微信通知
func sendWecom(message string) {
	webhook := config.AppConfig.Notification.Wecom.Webhook
	if webhook == "" {
		return
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send wecom notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("Wecom notification sent successfully")
	} else {
		log.Printf("Wecom notification failed with status: %d", resp.StatusCode)
	}
}

// sendDingtalk 发送钉钉通知
func sendDingtalk(message string) {
	webhook := config.AppConfig.Notification.Dingtalk.Webhook
	secret := config.AppConfig.Notification.Dingtalk.Secret

	if webhook == "" {
		return
	}

	// 如果配置了加签，计算签名
	if secret != "" {
		timestamp := time.Now().UnixMilli()
		sign := calculateDingtalkSign(timestamp, secret)
		webhook = fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, sign)
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send dingtalk notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("Dingtalk notification sent successfully")
	} else {
		log.Printf("Dingtalk notification failed with status: %d", resp.StatusCode)
	}
}

// calculateDingtalkSign 计算钉钉加签
func calculateDingtalkSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
