package main

import "fmt"

func getLangPrompt(lang string) string {
	switch lang {
	case "en":
		return "领导是外国人，请用国外的书信格式来写，全文要用英语"
	default:
		return "领导是中国人，请用中国的书信格式来写"
	}
}

func getDayPrompt(lang string) string {
	return fmt.Sprintf(`我是一名程序员，我想用git日志来向领导表达我今天的工作内容，帮我根据git日志生成我今天的工作内容报道。
    要求你以我的身份第一视角来写这份报告，不要提到是你帮助我写的，也不要透露是根据git日志生成的。
    不要提到commit之类的git概念
    需要分点来写
    %s
    我的git日志如下:`, getLangPrompt(lang))
}

func getWeekPrompt(lang string) string {
	return fmt.Sprintf(`我是一名程序员，我想用git日志来向领导表达我本周的工作内容，帮我根据git日志生成我本周的工作内容报道。
    要求你以我的身份第一视角来写这份报告，不要提到是你帮助我写的，不要透露是根据git日志生成的。
    不要提到commit之类的git概念
    需要分点来写
    需要展现出我对之前工作内容的总结，以及对下周工作的展望
	%s
    我的git日志如下:`, getLangPrompt(lang))
}
