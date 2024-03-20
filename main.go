package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os/exec"
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
	gitLog := gitShow()
	aiReq(
		getAiPrompt(),
		&gitLog,
	)
}

func execCmd(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func getEndGitUsername() string {
	cutset := " \t\r\n"
	username := strings.Trim(getGitUsername(), cutset)
	if len(username) == 0 {
		return execCmd("git", "config", "--global", "user.name")
	}
	return username
}

func gitShow() string {
	// 获取当前日期-1
	now := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
	afterCmd := "--after=" + now

	committerCmd := "--committer=" + getEndGitUsername()

	dirCmd := "--git-dir=" + getGitPath()

	switch getReportMode() {
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

func aiReq(prompt string, gitLog *string) {
	gitLogText := *gitLog
	log.Debug("gitLog length:", len(gitLogText))
	if len(gitLogText) >= 7800 {
		gitLogText = gitLogText[:7800]
	}

	client := resty.New()

	log.Info("请求中......")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("Bearer " + getAiAk()).
		SetBody(ApiReqBody{
			Model: getAiModel(),
			Input: ApiReqBodyInput{
				// TODO: 根据model定义token最大长度
				Prompt: prompt + gitLogText,
			},
		}).
		Post(getAiUrl())

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
