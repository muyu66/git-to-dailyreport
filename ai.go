package main

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type Ai interface {
	request(*resty.Client, string, *string) string
}

type AliyunAi struct {
}

type OllamaAi struct {
}

func AiFactory(aiName string) (Ai, error) {
	switch aiName {
	case "aliyun":
		return AliyunAi{}, nil
	case "ollama":
		return OllamaAi{}, nil
	default:
		return nil, errors.New("找不到合适的AI引擎")
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

type OllamaApiReqBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaApiRes struct {
	Response string `json:"response"`
}

func (ai AliyunAi) request(client *resty.Client, prompt string, gitLog *string) string {
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
		Post("https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation")

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

func (ai OllamaAi) request(client *resty.Client, prompt string, gitLog *string) string {
	log.Info("请求大模型中......")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("Bearer " + getAiAkConf()).
		SetBody(OllamaApiReqBody{
			Model:  getAiModelConf(),
			Prompt: prompt + *gitLog,
			Stream: false,
		}).
		SetDoNotParseResponse(true).
		Post("http://localhost:11434/api/generate")

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != 200 {
		log.Fatal(resp)
	}

	log.Info("请求完成......")

	result := OllamaApiRes{}
	jsonErr := json.Unmarshal(resp.Body(), &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result.Response
}
