package main

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func flow(ai Ai, client *resty.Client, prompt string) AiReqBodyMessage {
	var flowPrompts = getFlowPrompt()

	log.Info("开启AI-FLOW......")
	var userMsgs = []AiReqBodyMessage{
		{
			Role:    "system",
			Content: flowPrompts[0],
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
	var msgs1 = ai.request(client, userMsgs)

	var messagesBoss = []AiReqBodyMessage{
		{
			Role:    "system",
			Content: flowPrompts[1],
		},
		{
			Role:    "user",
			Content: msgs1.Content,
		},
	}
	var msgs2 = ai.request(client, messagesBoss)

	userMsgs = append(userMsgs, msgs1)
	userMsgs = append(userMsgs, AiReqBodyMessage{
		Role:    "user",
		Content: flowPrompts[2] + msgs2.Content,
	})
	var res = ai.request(client, userMsgs)
	log.Info("已关闭AI-FLOW......")
	return res
}
