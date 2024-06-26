package main

import "fmt"

func getLangPrompt(lang string) string {
	switch lang {
	case "en":
		return "领导是外国人，请用国外的汇报格式来写，全文要用英语"
	default:
		return "领导是中国人，请用中国的汇报格式来写"
	}
}

func basePrompt() string {
	return `我是一名程序员，因为GIT里面存储着我的工作代码，所以我想用抽取到的GIT日志来向领导表达我今天的工作内容，
	帮我根据git日志生成我今天的工作内容报道。
    要求你以我的身份第一视角来写这份报告，不要提到是你帮助我写的，也不要透露是根据git日志生成的。
    不要提到commit之类的git概念
    需要分点来写`
}

func getDayPrompt(lang string) string {
	return fmt.Sprintf(`
	%s
    %s
    我的git日志如下:`, basePrompt(), getLangPrompt(lang))
}

func getWeekPrompt(lang string) string {
	return fmt.Sprintf(`
	%s
    需要展现出我对之前工作内容的总结，以及对下周工作的展望
	%s
    我的git日志如下:`, basePrompt(), getLangPrompt(lang))
}

func getFlowPrompt() [3]string {
	return [3]string{
		"你是一个擅于帮助程序员分析GIT日志，并整理归纳今天工作内容的助理。",
		"你是一个程序员的领导，你要审查他提交的今日工作报告，并提出修改意见。",
		"我将你的报告提交给了我的领导，他反馈了一些修改建议，请你进行补充，并重新提交。建议如下：",
	}
}
