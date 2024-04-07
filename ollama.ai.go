package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type OllamaAi struct {
	Config AiConfig
}

type OllamaAiReqBody struct {
	Model    string             `json:"model"`
	Messages []AiReqBodyMessage `json:"messages"`
	Stream   bool               `json:"stream"`
}

type OllamaAiRes struct {
	Message AiReqBodyMessage `json:"message"`
}

func (ai OllamaAi) before() {
	log.Info("请求大模型中......")
}

func (ai OllamaAi) after() {
	log.Info("请求完成......")
}

func (ai OllamaAi) request(client *resty.Client, messages []AiReqBodyMessage) AiReqBodyMessage {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(OllamaAiReqBody{
			Model:    ai.Config.Model,
			Messages: messages,
			Stream:   false,
		}).
		Post(ai.Config.BaseUrl + "/api/chat")

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != 200 {
		log.Fatal(resp)
	}

	result := OllamaAiRes{}
	jsonErr := json.Unmarshal(resp.Body(), &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result.Message
}
