package scanner

import (
	"bufio"
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
	Timestamp     string   `json:"timestamp"`
	URL           string   `json:"url"`
	Input         string   `json:"input"`
	Title         string   `json:"title"`
	StatusCode    int      `json:"status_code"`
	ContentLength int      `json:"content_length"`
	Scheme        string   `json:"scheme"`
	Webserver     string   `json:"webserver"`
	ContentType   string   `json:"content_type"`
	Method        string   `json:"method"`
	Host          string   `json:"host"`
	Path          string   `json:"path"`
	Time          string   `json:"time"`
	A             []string `json:"a"`
	Tech          []string `json:"tech"`
	CDN           bool     `json:"cdn"`
	CDNName       string   `json:"cdn_name"`
	Failed        bool     `json:"failed"`
}

// ScanTarget 扫描指定目标
func ScanTarget(targetID uint, domain string) error {
	log.Printf("Starting scan for target: %s (ID: %d)", domain, targetID)

	// 创建扫描任务
	task := &models.ScanTask{
		TargetID:    targetID,
		Domain:      domain,
		TaskType:    "subdomain",
		Status:      "running",
		Progress:    0,
		CurrentStep: "初始化",
	}
	now := time.Now()
	task.StartedAt = &now
	database.DB.Create(task)

	// 创建日志记录器
	logger := NewLogger(task.ID)
	return executeScan(task.ID, targetID, domain, logger)
}

// ScanWithTask 使用已存在的任务进行扫描
func ScanWithTask(taskID uint, targetID uint, domain string) error {
	log.Printf("Starting scan with task ID: %d for domain: %s", taskID, domain)

	// 更新任务状态
	now := time.Now()
	database.DB.Model(&models.ScanTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":       "running",
		"started_at":   &now,
		"progress":     0,
		"current_step": "初始化",
	})

	// 创建日志记录器
	logger := NewLogger(taskID)
	return executeScan(taskID, targetID, domain, logger)
}

// executeScan 执行扫描的核心逻辑
func executeScan(taskID uint, targetID uint, domain string, logger *Logger) error {
	logger.Info("初始化", fmt.Sprintf("开始扫描目标: %s", domain))
	logger.UpdateProgress(5, "子域名收集")

	// 1. 子域名收集
	logger.Info("子域名收集", "开始收集子域名...")
	subdomains, err := collectSubdomains(domain, logger)
	if err != nil {
		logger.Error("子域名收集", fmt.Sprintf("收集失败: %v", err))
		updateTaskFailedByID(taskID, err)
		return err
	}
	logger.Success("子域名收集", fmt.Sprintf("收集到 %d 个子域名", len(subdomains)))
	logger.UpdateProgress(40, "存活验证")

	// 2. 存活验证
	logger.Info("存活验证", fmt.Sprintf("开始验证 %d 个子域名...", len(subdomains)))
	aliveSubdomains, err := verifyAlive(subdomains, logger)
	if err != nil {
		logger.Error("存活验证", fmt.Sprintf("验证失败: %v", err))
		updateTaskFailedByID(taskID, err)
		return err
	}
	logger.Success("存活验证", fmt.Sprintf("发现 %d 个存活子域名", len(aliveSubdomains)))
	logger.UpdateProgress(70, "差量对比")

	// 3. 差量对比
	logger.Info("差量对比", "开始对比历史数据...")
	changes := compareAndUpdate(targetID, aliveSubdomains, logger)
	logger.Success("差量对比", fmt.Sprintf("发现 %d 条变更", len(changes)))
	logger.UpdateProgress(90, "发送通知")

	// 4. 发送通知
	if len(changes) > 0 {
		logger.Info("发送通知", fmt.Sprintf("发送 %d 条变更通知...", len(changes)))
		notifier.SendChangeNotification(domain, changes)
		logger.Success("发送通知", "通知发送成功")
	} else {
		logger.Info("发送通知", "无变更，跳过通知")
	}

	// 更新任务状态
	logger.UpdateProgress(100, "完成")
	updateTaskCompletedByID(taskID, fmt.Sprintf("发现 %d 个存活子域名，%d 条变更", len(aliveSubdomains), len(changes)))
	logger.Success("完成", fmt.Sprintf("扫描完成: %s", domain))

	log.Printf("Scan completed for %s", domain)
	return nil
}

// collectSubdomains 收集子域名
func collectSubdomains(domain string, logger *Logger) ([]string, error) {
	dataDir := config.AppConfig.Scanner.DataDir
	subdomainSet := make(map[string]bool)

	// 1. Subfinder
	logger.Info("Subfinder", "运行 Subfinder...")
	log.Printf("Running subfinder for %s", domain)

	subfinderOutput := filepath.Join(dataDir, fmt.Sprintf("%s_subfinder.txt", domain))
	cmd := exec.Command("subfinder", "-d", domain, "-all", "-silent", "-o", subfinderOutput)

	if err := cmd.Run(); err != nil {
		logger.Warning("Subfinder", fmt.Sprintf("Subfinder 执行出错: %v", err))
		log.Printf("Subfinder error: %v", err)
	} else {
		// 读取结果
		count := readSubdomainsFromFile(subfinderOutput, subdomainSet)
		logger.Success("Subfinder", fmt.Sprintf("发现 %d 个子域名", count))
	}

	// 2. Samoscout
	logger.Info("Samoscout", "运行 Samoscout...")
	log.Printf("Running samoscout for %s", domain)

	samoscoutOutput := filepath.Join(dataDir, fmt.Sprintf("%s_samoscout.txt", domain))
	cmd = exec.Command("samoscout", "-d", domain, "-silent", "-o", samoscoutOutput)

	if err := cmd.Run(); err != nil {
		logger.Warning("Samoscout", fmt.Sprintf("Samoscout 执行出错: %v", err))
		log.Printf("Samoscout error: %v", err)
	} else {
		count := readSubdomainsFromFile(samoscoutOutput, subdomainSet)
		logger.Success("Samoscout", fmt.Sprintf("发现 %d 个子域名", count))
	}

	// 转换为切片
	subdomains := make([]string, 0, len(subdomainSet))
	for sub := range subdomainSet {
		subdomains = append(subdomains, sub)
	}

	// 保存所有子域名到文件
	allSubsFile := filepath.Join(dataDir, fmt.Sprintf("%s_all_subs.txt", domain))
	file, err := os.Create(allSubsFile)
	if err != nil {
		return subdomains, err
	}
	defer file.Close()

	for _, sub := range subdomains {
		file.WriteString(sub + "\n")
	}

	logger.Info("汇总", fmt.Sprintf("共收集到 %d 个唯一子域名", len(subdomains)))

	// 4. Ksubdomain 验证 DNS 存活
	if len(subdomains) > 0 {
		logger.Info("Ksubdomain", fmt.Sprintf("验证 %d 个子域名的 DNS 存活...", len(subdomains)))
		verifiedFile := filepath.Join(dataDir, fmt.Sprintf("%s_verified.txt", domain))

		cmd = exec.Command("ksubdomain", "verify", "-f", allSubsFile, "--silent", "-o", verifiedFile)
		if err := cmd.Run(); err != nil {
			logger.Warning("Ksubdomain", fmt.Sprintf("Ksubdomain 执行出错: %v", err))
		} else {
			// 读取验证后的子域名
			verifiedSet := make(map[string]bool)
			count := readSubdomainsFromFile(verifiedFile, verifiedSet)

			// 更新子域名列表为验证后的
			if count > 0 {
				subdomains = make([]string, 0, len(verifiedSet))
				for sub := range verifiedSet {
					subdomains = append(subdomains, sub)
				}
				logger.Success("Ksubdomain", fmt.Sprintf("验证通过 %d 个存活子域名", count))
			}
		}
	}

	return subdomains, nil
}

// readSubdomainsFromFile 从文件读取子域名到 set
func readSubdomainsFromFile(filename string, subdomainSet map[string]bool) int {
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sub := strings.TrimSpace(scanner.Text())
		if sub != "" && !subdomainSet[sub] {
			subdomainSet[sub] = true
			count++
		}
	}
	return count
}

// queryCertSh 查询 cert.sh
// verifyAlive 验证存活
func verifyAlive(subdomains []string, logger *Logger) ([]HttpxResult, error) {
	if len(subdomains) == 0 {
		logger.Warning("Httpx", "没有子域名需要验证")
		return []HttpxResult{}, nil
	}

	dataDir := config.AppConfig.Scanner.DataDir
	inputFile := filepath.Join(dataDir, "temp_subs.txt")
	outputFile := filepath.Join(dataDir, "temp_httpx.json")

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
	logger.Info("Httpx", fmt.Sprintf("开始验证 %d 个子域名...", len(subdomains)))
	log.Printf("Running httpx to verify %d subdomains", len(subdomains))

	cmd := exec.Command("httpx",
		"-l", inputFile,
		"-sc",     // 状态码
		"-title",  // 标题
		"-td",     // 技术栈
		"-json",   // JSON 输出
		"-silent", // 静默模式
		"-o", outputFile,
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("httpx error: %w", err)
	}

	// 解析结果
	var results []HttpxResult
	resultFile, err := os.Open(outputFile)
	if err != nil {
		return nil, err
	}
	defer resultFile.Close()

	scanner := bufio.NewScanner(resultFile)
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

	// 清理临时文件
	os.Remove(inputFile)

	logger.Success("Httpx", fmt.Sprintf("验证完成，发现 %d 个存活子域名", len(results)))

	return results, nil
}

// compareAndUpdate 对比并更新数据库
func compareAndUpdate(targetID uint, aliveSubdomains []HttpxResult, logger *Logger) []models.ChangeLog {
	var changes []models.ChangeLog
	now := time.Now()

	logger.Info("对比", "获取历史数据...")

	// 获取数据库中已有的子域名
	var existingSubdomains []models.Subdomain
	database.DB.Where("target_id = ?", targetID).Find(&existingSubdomains)

	existingMap := make(map[string]*models.Subdomain)
	for i := range existingSubdomains {
		existingMap[existingSubdomains[i].Subdomain] = &existingSubdomains[i]
	}

	logger.Info("对比", fmt.Sprintf("开始对比 %d 个新结果与 %d 条历史记录...", len(aliveSubdomains), len(existingSubdomains)))

	newCount := 0
	aliveCount := 0
	deadCount := 0

	// 处理新扫描到的存活子域名
	for _, alive := range aliveSubdomains {
		subdomain := extractDomain(alive.URL)

		// 提取技术栈
		tech := ""
		if len(alive.Tech) > 0 {
			tech = strings.Join(alive.Tech, ", ")
		}

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
				aliveCount++
			}

			// 更新记录
			existing.Status = "alive"
			existing.StatusCode = alive.StatusCode
			existing.Title = alive.Title
			existing.Server = alive.Webserver
			existing.IP = alive.Host
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
				Server:     alive.Webserver,
				IP:         alive.Host,
				FirstSeen:  now,
				LastSeen:   now,
			}
			database.DB.Create(&newSub)

			change := models.ChangeLog{
				TargetID:   targetID,
				ChangeType: "subdomain_new",
				Content:    fmt.Sprintf("%s (状态码: %d, 标题: %s, 技术: %s)", subdomain, alive.StatusCode, alive.Title, tech),
			}
			changes = append(changes, change)
			database.DB.Create(&change)
			newCount++
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
			deadCount++
		}
	}

	logger.Info("对比", fmt.Sprintf("对比完成 - 新增: %d, 恢复: %d, 失效: %d", newCount, aliveCount, deadCount))

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

func updateTaskFailedByID(taskID uint, err error) {
	now := time.Now()
	database.DB.Model(&models.ScanTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":       "failed",
		"error_msg":    err.Error(),
		"completed_at": &now,
	})
}

func updateTaskCompleted(task *models.ScanTask, result string) {
	task.Status = "completed"
	task.Result = result
	now := time.Now()
	task.CompletedAt = &now
	database.DB.Save(task)
}

func updateTaskCompletedByID(taskID uint, result string) {
	now := time.Now()
	database.DB.Model(&models.ScanTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":       "completed",
		"result":       result,
		"completed_at": &now,
	})
}
