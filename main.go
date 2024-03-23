package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var reportCycle string
var reportModeConf string
var reportLangConf string

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// 初始化配置管理器
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Info("Error reading config file:", err)
	} else {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

	// 初始化全局配置
	reportModeConf = getReportModeConf()
	reportLangConf = getReportLangConf()
}

func loadCmdParams() {
	// 从命令行参数读取配置
	flag.StringVar(&reportCycle, "c", "week", "报告周期[day|week]")
	flag.Parse()
	log.Debug("[命令行参数] reportCycle=", reportCycle)
}

// TODO: 调整代码层次结构
func main() {
	loadCmdParams()
	prompt := getAiPrompt(reportCycle, reportLangConf)

	gitLog := fire(prompt)

	client := resty.New()
	reportText := aiReq(
		client,
		prompt,
		&gitLog,
	)
	out(&reportText)
}

func formatText(content string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(content, "")
}

func getAiPrompt(dayOrWeek string, lang string) string {
	if dayOrWeek == "week" {
		return getWeekPrompt(lang)
	}
	return getDayPrompt(lang)
}

func out(reportText *string) {
	switch getReportOutConf() {
	case "file":
		writeFile(reportText)
	case "print":
		fmt.Println(*reportText)
	default:
		log.Fatal("不支持的输出方式：", getReportOutConf())
	}
}

func writeFile(reportText *string) {
	nowDate := time.Now().Format(time.DateOnly)
	file, err := os.Create(nowDate + ".txt")
	if err != nil {
		log.Fatal("An error occurred while opening the file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("An error occurred while opening the file:", err)
		}
	}(file)

	// 将字符串写入文件
	_, err = io.WriteString(file, *reportText)
	if err != nil {
		log.Fatal("An error occurred while writing to file:", err)
		return
	}
}

func getCmdInfoDate(dayOrWeek string) string {
	const timeStr = " 00:00:00"
	if dayOrWeek == "week" {
		return time.Now().AddDate(0, 0, -7+1).Format(time.DateOnly) + timeStr
	}
	return time.Now().AddDate(0, 0, -getReportIntervalDayConf()+1).Format(time.DateOnly) + timeStr
}

func fire(prompt string) string {
	var maxLogLen int64 = 0
	var gitLogMap = sync.Map{}

	var cmdInfo = CmdInfo{
		Username: getEndGitUsername(),
		Date:     getCmdInfoDate(reportCycle),
	}

	var wg sync.WaitGroup
	for _, repo := range getGitReposConf() {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			getGitLogs(repo, &gitLogMap, &maxLogLen, cmdInfo)
		}(repo)
	}
	wg.Wait()

	gitLog := makeAiReq(&gitLogMap, maxLogLen, prompt)
	log.Debug("maxLogLen=", maxLogLen)
	log.Debug("len(gitLog)=", len(gitLog))
	return gitLog
}

func makeAiReq(gitLogMap *sync.Map, maxLogLen int64, prompt string) string {
	if maxLogLen <= 0 {
		log.Fatal("未发现存在日志内容，请检查日志，退出")
	}

	gitLogMapKeyCount := 0
	gitLogMap.Range(func(_, _ interface{}) bool {
		gitLogMapKeyCount++
		return true
	})

	gitLogLarge := ""
	aiModel := getAiModel(getAiModelConf())
	// 大模型最多支持一次传输TOKEN数量
	maxInputTokenCount := int(aiModel.MaxInputTokenCount*1000) - len(prompt) - (gitLogMapKeyCount * 15) - 10
	// 每个GIT仓库可以分配到多少TOKEN的比例
	ratio := float64(maxInputTokenCount) / float64(maxLogLen)
	log.Debug("maxInputTokenCount=", maxInputTokenCount)
	log.Debug("maxLogLen=", maxLogLen)
	log.Debug("ratio=", ratio)

	gitLogMap.Range(func(key, value interface{}) bool {
		repoName := key.(string)
		gitLog := value.(string)
		maxGitLogLen := int(float64(len(gitLog)) * ratio)
		log.Debug("repoName=", repoName)
		log.Debug("gitLog原始长度=", len(gitLog))
		log.Debug("gitLog最大长度=", maxGitLogLen)
		if len(gitLog) >= maxGitLogLen {
			gitLog = gitLog[:maxGitLogLen]
		}
		log.Debug("gitLog最终长度=", len(gitLog))
		gitLogLarge += "\n以下是" + repoName + "仓库的GIT日志\n" + gitLog
		return true
	})
	return gitLogLarge
}

func getGitLogs(repo string, gitLogMap *sync.Map, maxLogLen *int64, cmdInfo CmdInfo) {
	// 调用CMD获取指定目录下的所有GIT仓库地址
	// TODO: 改成内置方法，以便于支撑跨平台
	repoPathStr, err := execCmd("cmd", "/c", "cd "+repo+" && for /d /r %i in (.git) do @if exist %i\\HEAD echo %i")
	if err != nil {
		log.Fatal("获取GIT仓库地址异常")
	}
	// 整理换行符
	repoPathStr = strings.Replace(repoPathStr, "\r\n", "\n", -1)
	repoPaths := strings.Split(repoPathStr, "\n")

	var wg sync.WaitGroup
	for _, repoPath := range repoPaths {
		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()
			if len(repoPath) == 0 {
				return
			}
			gitLog, err := getGitLog(repoPath, cmdInfo)
			if err != nil {
				log.Error("[异常] 获取工作日志异常", repoPath)
			}
			// 压缩文本体积
			gitLog = formatText(gitLog)
			if len(gitLog) == 0 {
				log.Info("未找到工作日志", repoPath)
			} else {
				getGitLogAfter(repoPath, &gitLog, gitLogMap, maxLogLen)
				log.Info("已加载", repoPath)
			}
		}(repoPath)
	}
	wg.Wait()
}

func execCmd(name string, arg ...string) (string, error) {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type CmdInfo struct {
	Username string
	Date     string
}

func getEndGitUsername() string {
	cutset := " \t\r\n"
	username := strings.Trim(getGitUsernameConf(), cutset)
	if len(username) == 0 {
		res, err := execCmd("git", "config", "--global", "user.name")
		if err != nil {
			log.Fatal("获取GIT用户名异常")
		}
		return res
	}
	return username
}

func getGitLogAfter(repoPath string, gitLog *string, gitLogMap *sync.Map, maxLogLen *int64) {
	repoName := filepath.Base(filepath.Dir(repoPath))
	gitLogMap.Store(repoName, *gitLog)
	// 后续需要用到repoPath，于是提前加入，防止长度溢出
	atomic.AddInt64(maxLogLen, int64(len(*gitLog)+len(repoName)))
}

func getGitLog(repoPath string, cmdInfo CmdInfo) (string, error) {
	// 获取当前日期-1
	afterCmd := "--after=\"" + cmdInfo.Date + "\""

	committerCmd := "--committer=" + cmdInfo.Username

	dirCmd := "--git-dir=" + repoPath

	switch reportModeConf {
	case "normal":
		return execCmd("git", dirCmd, "log", "--stat", "-p", afterCmd, committerCmd)
	case "simple":
		return execCmd("git", dirCmd, "log", "--stat", afterCmd, committerCmd)
	default:
		return execCmd("git", dirCmd, "log", "--stat", afterCmd, committerCmd)
	}
}

type ApiReqBody struct {
	Model string          `json:"model"`
	Input ApiReqBodyInput `json:"input"`
	//Messages []ApiReqBodyMessage `json:"messages"`
}

type ApiReqBodyMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiReqBodyInput struct {
	Prompt string `json:"prompt"`
}

type ApiRes struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
	} `json:"output"`
	RequestID string `json:"request_id"`
}

func aiReq(client *resty.Client, prompt string, gitLog *string) string {
	log.Info("请求大模型中......")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("Bearer " + getAiAkConf()).
		SetBody(ApiReqBody{
			Model: getAiModelConf(),
			Input: ApiReqBodyInput{
				Prompt: prompt + *gitLog,
			},
		}).
		Post(getAiUrlConf())

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != 200 {
		log.Fatal(resp)
	}

	log.Info("请求完成......")

	result := ApiRes{}
	jsonErr := json.Unmarshal(resp.Body(), &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result.Output.Text
}
