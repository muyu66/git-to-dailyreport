package main

import "github.com/go-resty/resty/v2"

type OpenAi struct {
}

// TODO:
func (ai OpenAi) request(client *resty.Client, messages []AiReqBodyMessage) AiReqBodyMessage {
	return AiReqBodyMessage{}
}
