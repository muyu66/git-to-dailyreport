package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type OpenAi struct {
	Config AiConfig
}

type OpenAiReqBody struct {
	Model       string             `json:"model"`
	Messages    []AiReqBodyMessage `json:"messages"`
	Temperature float64            `json:"temperature"`
}

type OpenAiRes struct {
	Choices []AiResChoice `json:"choices"`
}

func (ai OpenAi) before() {
	log.Info("请求大模型中......")
}

func (ai OpenAi) after() {
	log.Info("请求完成......")
}

func (ai OpenAi) request(client *resty.Client, messages []AiReqBodyMessage) AiReqBodyMessage {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(ai.Config.ApiKey).
		SetBody(OpenAiReqBody{
			Model:       ai.Config.Model,
			Messages:    messages,
			Temperature: 0.7,
		}).
		Post(ai.Config.BaseUrl + "/v1/chat/completions")

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != 200 {
		log.Fatal(resp)
	}

	result := OpenAiRes{}
	jsonErr := json.Unmarshal(resp.Body(), &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result.Choices[0].Message
}
