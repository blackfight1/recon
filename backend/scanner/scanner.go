package scanner

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"recon-platform/config"
	"recon-platform/database"
	"recon-platform/models"
	"recon-platform/notifier"
	"strings"
	"time"
)

// HttpxResult Httpx 输出结构
type HttpxResult struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Title      string `json:"title"`
	Server     string `json:"webserver"`
	IP         string `json:"host"`
}

// ScanTarget 扫描指定目标
func ScanTarget(targetID uint, domain string) error {
	log.Printf("Starting scan for target: %s (ID: %d)", domain, targetID)

	// 创建扫描任务
	task := &models.ScanTask{
		TargetID: targetID,
		TaskType: "subdomain",
		Status:   "running",
	}
	now := time.Now()
	task.StartedAt = &now
	database.DB.Create(task)

	// 1. 子域名收集
	subdomains, err := collectSubdomains(domain)
	if err != nil {
		updateTaskFailed(task, err)
		return err
	}
	log.Printf("Collected %d subdomains for %s", len(subdomains), domain)

	// 2. 存活验证
	aliveSubdomains, err := verifyAlive(subdomains)
	if err != nil {
		updateTaskFailed(task, err)
		return err
	}
	log.Printf("Found %d alive subdomains for %s", len(aliveSubdomains), domain)

	// 3. 差量对比
	changes := compareAndUpdate(targetID, aliveSubdomains)

	// 4. 发送通知
	if len(changes) > 0 {
		notifier.SendChangeNotification(domain, changes)
	}

	// 更新任务状态
	updateTaskCompleted(task, fmt.Sprintf("Found %d alive subdomains, %d changes", len(aliveSubdomains), len(changes)))

	log.Printf("Scan completed for %s", domain)
	return nil
}

// collectSubdomains 收集子域名
func collectSubdomains(domain string) ([]string, error) {
	dataDir := config.AppConfig.Scanner.DataDir
	outputFile := filepath.Join(dataDir, fmt.Sprintf("%s_subs.txt", domain))

	subdomainSet := make(map[string]bool)

	// 1. Subfinder
	log.Printf("Running subfinder for %s", domain)
	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/output", dataDir),
		"projectdiscovery/subfinder:latest",
		"-d", domain,
		"-all",
		"-silent",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Subfinder error: %v", err)
	} else {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			sub := strings.TrimSpace(scanner.Text())
			if sub != "" {
				subdomainSet[sub] = true
			}
		}
	}

	// 2. Assetfinder
	log.Printf("Running assetfinder for %s", domain)
	cmd = exec.Command("docker", "run", "--rm",
		"tomnomnom/assetfinder:latest",
		"--subs-only",
		domain,
	)
	output, err = cmd.Output()
	if err != nil {
		log.Printf("Assetfinder error: %v", err)
	} else {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			sub := strings.TrimSpace(scanner.Text())
			if sub != "" {
				subdomainSet[sub] = true
			}
		}
	}

	// 3. cert.sh (被动查询)
	log.Printf("Querying cert.sh for %s", domain)
	certSubs, err := queryCertSh(domain)
	if err != nil {
		log.Printf("cert.sh error: %v", err)
	} else {
		for _, sub := range certSubs {
			subdomainSet[sub] = true
		}
	}

	// 转换为切片
	subdomains := make([]string, 0, len(subdomainSet))
	for sub := range subdomainSet {
		subdomains = append(subdomains, sub)
	}

	// 保存到文件
	file, err := os.Create(outputFile)
	if err != nil {
		return subdomains, err
	}
	defer file.Close()

	for _, sub := range subdomains {
		file.WriteString(sub + "\n")
	}

	return subdomains, nil
}

// queryCertSh 查询 cert.sh
func queryCertSh(domain string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "curl", "-s",
		fmt.Sprintf("https://crt.sh/?q=%%.%s&output=json", domain))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var results []struct {
		NameValue string `json:"name_value"`
	}

	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	subdomainSet := make(map[string]bool)
	for _, result := range results {
		subs := strings.Split(result.NameValue, "\n")
		for _, sub := range subs {
			sub = strings.TrimSpace(sub)
			if sub != "" && !strings.HasPrefix(sub, "*") {
				subdomainSet[sub] = true
			}
		}
	}

	subdomains := make([]string, 0, len(subdomainSet))
	for sub := range subdomainSet {
		subdomains = append(subdomains, sub)
	}

	return subdomains, nil
}

// verifyAlive 验证存活
func verifyAlive(subdomains []string) ([]HttpxResult, error) {
	if len(subdomains) == 0 {
		return []HttpxResult{}, nil
	}

	dataDir := config.AppConfig.Scanner.DataDir
	inputFile := filepath.Join(dataDir, "temp_subs.txt")
	outputFile := filepath.Join(dataDir, "temp_alive.json")

	// 写入临时文件
	file, err := os.Create(inputFile)
	if err != nil {
		return nil, err
	}
	for _, sub := range subdomains {
		file.WriteString(sub + "\n")
	}
	file.Close()

	// 运行 httpx
	log.Printf("Running httpx to verify %d subdomains", len(subdomains))
	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/output", dataDir),
		"projectdiscovery/httpx:latest",
		"-l", "/output/temp_subs.txt",
		"-json",
		"-silent",
		"-timeout", "10",
		"-retries", "2",
		"-threads", "50",
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("httpx error: %w", err)
	}

	// 解析结果
	var results []HttpxResult
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var result HttpxResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			log.Printf("Failed to parse httpx result: %s", line)
			continue
		}
		results = append(results, result)
	}

	// 保存结果
	os.WriteFile(outputFile, output, 0644)

	// 清理临时文件
	os.Remove(inputFile)

	return results, nil
}

// compareAndUpdate 对比并更新数据库
func compareAndUpdate(targetID uint, aliveSubdomains []HttpxResult) []models.ChangeLog {
	var changes []models.ChangeLog
	now := time.Now()

	// 获取数据库中已有的子域名
	var existingSubdomains []models.Subdomain
	database.DB.Where("target_id = ?", targetID).Find(&existingSubdomains)

	existingMap := make(map[string]*models.Subdomain)
	for i := range existingSubdomains {
		existingMap[existingSubdomains[i].Subdomain] = &existingSubdomains[i]
	}

	// 处理新扫描到的存活子域名
	for _, alive := range aliveSubdomains {
		subdomain := extractDomain(alive.URL)

		if existing, found := existingMap[subdomain]; found {
			// 已存在，更新状态
			if existing.Status != "alive" {
				// 状态变化：从 dead 变为 alive
				change := models.ChangeLog{
					TargetID:   targetID,
					ChangeType: "subdomain_alive",
					Content:    fmt.Sprintf("%s (状态码: %d, 标题: %s)", subdomain, alive.StatusCode, alive.Title),
				}
				changes = append(changes, change)
				database.DB.Create(&change)
			}

			// 更新记录
			existing.Status = "alive"
			existing.StatusCode = alive.StatusCode
			existing.Title = alive.Title
			existing.Server = alive.Server
			existing.IP = alive.IP
			existing.LastSeen = now
			database.DB.Save(existing)

			delete(existingMap, subdomain)
		} else {
			// 新发现的子域名
			newSub := models.Subdomain{
				TargetID:   targetID,
				Subdomain:  subdomain,
				Status:     "new",
				StatusCode: alive.StatusCode,
				Title:      alive.Title,
				Server:     alive.Server,
				IP:         alive.IP,
				FirstSeen:  now,
				LastSeen:   now,
			}
			database.DB.Create(&newSub)

			change := models.ChangeLog{
				TargetID:   targetID,
				ChangeType: "subdomain_new",
				Content:    fmt.Sprintf("%s (状态码: %d, 标题: %s)", subdomain, alive.StatusCode, alive.Title),
			}
			changes = append(changes, change)
			database.DB.Create(&change)
		}
	}

	// 处理消失的子域名（超过 24 小时未见）
	for _, existing := range existingMap {
		if existing.Status == "alive" && time.Since(existing.LastSeen) > 24*time.Hour {
			existing.Status = "dead"
			database.DB.Save(existing)

			change := models.ChangeLog{
				TargetID:   targetID,
				ChangeType: "subdomain_dead",
				Content:    existing.Subdomain,
			}
			changes = append(changes, change)
			database.DB.Create(&change)
		}
	}

	return changes
}

// extractDomain 从 URL 提取域名
func extractDomain(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	parts := strings.Split(url, "/")
	return parts[0]
}

func updateTaskFailed(task *models.ScanTask, err error) {
	task.Status = "failed"
	task.ErrorMsg = err.Error()
	now := time.Now()
	task.CompletedAt = &now
	database.DB.Save(task)
}

func updateTaskCompleted(task *models.ScanTask, result string) {
	task.Status = "completed"
	task.Result = result
	now := time.Now()
	task.CompletedAt = &now
	database.DB.Save(task)
}
