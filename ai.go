package main

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

type Ai interface {
	request(*resty.Client, []AiReqBodyMessage) AiReqBodyMessage
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

type AiReqBody struct {
	Model      string              `json:"model"`
	Input      AiReqBodyInput      `json:"input"`
	Parameters AiReqBodyParameters `json:"parameters"`
}

type AiReqBodyInput struct {
	Messages []AiReqBodyMessage `json:"messages"`
}

type AiReqBodyMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiResultFormat string

// 使用 iota 为枚举成员分配值
const (
	Message AiResultFormat = "message"
	Text    AiResultFormat = "text"
)

type AiReqBodyParameters struct {
	ResultFormat AiResultFormat `json:"result_format"`
}

type AiResChoice struct {
	Message AiReqBodyMessage `json:"message"`
	//FinishReason string           `json:"finish_reason"`
}

type AiRes struct {
	Output struct {
		Choices []AiResChoice `json:"choices"`
	} `json:"output"`
	RequestID string `json:"request_id"`
}
