package main

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func flow(ai Ai, client *resty.Client, prompt string) AiReqBodyMessage {
	log.Info("开始AI-FLOW......")
	var userMsgs = []AiReqBodyMessage{
		{
			Role:    "system",
			Content: "你是一个擅于帮助程序员分析GIT日志，并整理归纳今天工作内容的助理。",
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
			Content: "你是一个程序员的领导，你要审查他提交的今日工作报告，并提出修改意见。",
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
		Content: "我将你的报告提交给了我的领导，他反馈了一些修改建议，请你进行补充，并重新提交。建议如下：" + msgs2.Content,
	})
	var res = ai.request(client, userMsgs)
	log.Info("已结束AI-FLOW......")
	return res
}
