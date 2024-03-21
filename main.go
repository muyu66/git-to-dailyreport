package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Info("Error reading config file:", err)
	} else {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

	gitLogMap := make(map[string]string)
	maxLogLen := 0
	prompt := getAiPromptConf()

	for _, repo := range getGitReposConf() {
		getGitLogs(repo, &gitLogMap, &maxLogLen)
	}

	gitLog := makeAiReq(&gitLogMap, maxLogLen, prompt)

	client := resty.New()

	aiReq(
		client,
		getAiPromptConf(),
		&gitLog,
	)
}

func makeAiReq(gitLogMap *map[string]string, maxLogLen int, prompt string) string {
	gitLogLarge := ""
	gitLogM := *gitLogMap
	aiModel := getAiModel(getAiModelConf())
	// 大模型最多支持一次传输TOKEN数量
	maxInputTokenCount := int(aiModel.MaxInputTokenCount*1000) - len(prompt) - (len(*gitLogMap) * 15) - 10
	// 每个GIT仓库可以分配到多少TOKEN的比例
	ratio := float64(maxInputTokenCount) / float64(maxLogLen)
	log.Debug("maxInputTokenCount=", maxInputTokenCount)
	log.Debug("maxLogLen=", maxLogLen)
	log.Debug("ratio=", ratio)

	for repoName, gitLog := range gitLogM {
		maxGitLogLen := int(float64(len(gitLog)) * ratio)
		log.Debug("repoName=", repoName)
		log.Debug("gitLog原始长度=", len(gitLog))
		log.Debug("gitLog最大长度=", maxGitLogLen)
		if len(gitLog) >= maxGitLogLen {
			gitLog = gitLog[:maxGitLogLen]
		}
		log.Debug("gitLog最终长度=", len(gitLog))
		gitLogLarge += "\n以下是" + repoName + "仓库的GIT日志\n" + gitLog
	}

	return gitLogLarge
}

func getGitLogs(repo string, gitLogMap *map[string]string, maxLogLen *int) {
	gitLogM := *gitLogMap

	repoPathStr, err := execCmd("cmd", "/c", "cd "+repo+" && for /d /r %i in (.git) do @if exist %i\\HEAD echo %i")
	if err != nil {
		log.Fatal("获取GIT仓库地址异常")
	}
	// 整理换行符
	repoPathStr = strings.Replace(repoPathStr, "\r\n", "\n", -1)
	repoPaths := strings.Split(repoPathStr, "\n")

	for _, repoPath := range repoPaths {
		if len(repoPath) == 0 {
			continue
		}
		gitLog, err := getGitLog(repoPath)
		if err != nil {
			log.Error("[异常] 获取工作日志异常", repoPath)
		}
		if len(gitLog) == 0 {
			log.Info("未找到工作日志", repoPath)
		} else {
			repoName := filepath.Base(filepath.Dir(repoPath))
			gitLogM[repoName] = gitLog
			// 后续需要用到repoPath，于是提前加入，防止长度溢出
			*maxLogLen += len(gitLog) + len(repoName)
			log.Info("已加载", repoPath)
		}
	}
}

func execCmd(name string, arg ...string) (string, error) {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
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

func getGitLog(repoPath string) (string, error) {
	// 获取当前日期-1
	now := time.Now().AddDate(0, 0, -getReportIntervalDayConf()+1).Format(time.DateOnly)
	afterCmd := "--after=\"" + now + " 00:00:00\""

	committerCmd := "--committer=" + getEndGitUsername()

	dirCmd := "--git-dir=" + repoPath

	switch getReportModeConf() {
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

func aiReq(client *resty.Client, prompt string, gitLog *string) {
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
	fmt.Println(result.Output.Text)
}
