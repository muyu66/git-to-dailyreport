package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os/exec"
)

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

func gitShow() string {
	// 指定要执行的命令和参数
	cmd := exec.Command("git", "--git-dir="+getGitPath(), "show", getGitCommit())

	// 获取命令的标准输出
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
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
	client := resty.New()

	log.Info("请求中......")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("Bearer " + getAiAk()).
		SetBody(ApiReqBody{
			Model: getAiModel(),
			Input: ApiReqBodyInput{
				// TODO: 根据model定义token最大长度
				Prompt: prompt + gitLogText[:7800],
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
