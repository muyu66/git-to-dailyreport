package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type AliyunAi struct {
}

func (ai AliyunAi) request(client *resty.Client, messages []AiReqBodyMessage) AiReqBodyMessage {
	log.Info("请求大模型中......")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("Bearer " + getAiAkConf()).
		SetBody(AiReqBody{
			Model: getAiModelConf(),
			Input: AiReqBodyInput{
				Messages: messages,
			},
			Parameters: AiReqBodyParameters{
				ResultFormat: Message,
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

	result := AiRes{}
	jsonErr := json.Unmarshal(resp.Body(), &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result.Output.Choices[0].Message
}
