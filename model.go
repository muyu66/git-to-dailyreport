package main

type AiModel struct {
	Name               string
	MaxInputTokenCount uint16 // 单位K
}

var aiModels = []AiModel{
	// 默认
	// TODO: INPUT TOKEN考虑自定义覆盖
	{
		Name:               "other",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen-1.8b-chat",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen-7b-chat",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen-14b-chat",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen-72b-chat",
		MaxInputTokenCount: 30,
	},
	{
		Name:               "qwen1.5-7b-chat",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen1.5-14b-chat",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen1.5-72b-chat",
		MaxInputTokenCount: 30,
	},
	{
		Name:               "qwen-max-1201",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen-plus",
		MaxInputTokenCount: 30,
	},
	{
		Name:               "qwen-turbo",
		MaxInputTokenCount: 6,
	},
	{
		Name:               "qwen-max-longcontext",
		MaxInputTokenCount: 28,
	},
}

func getAiModel(name string) AiModel {
	for _, model := range aiModels {
		if name == model.Name {
			return model
		}
	}
	return aiModels[0]
}
