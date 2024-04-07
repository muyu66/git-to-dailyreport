package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type AliyunAi struct {
	Config AiConfig
}

func (ai AliyunAi) before() {
	log.Info("请求大模型中......")
}

func (ai AliyunAi) after() {
	log.Info("请求完成......")
}

func (ai AliyunAi) request(client *resty.Client, messages []AiReqBodyMessage) AiReqBodyMessage {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(ai.Config.ApiKey).
		SetBody(AiReqBody{
			Model: ai.Config.Model,
			Input: AiReqBodyInput{
				Messages: messages,
			},
			Parameters: AiReqBodyParameters{
				ResultFormat: Message,
			},
		}).
		Post(ai.Config.BaseUrl + "/api/v1/services/aigc/text-generation/generation")

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != 200 {
		log.Fatal(resp)
	}

	result := AiRes{}
	jsonErr := json.Unmarshal(resp.Body(), &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result.Output.Choices[0].Message
}
