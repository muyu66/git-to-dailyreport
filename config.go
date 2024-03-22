package main

import "github.com/spf13/viper"

func getAiPromptConf(dayOrWeek string) string {
	if dayOrWeek == "week" {
		return viper.GetString("aiPrompt.week")
	}
	return viper.GetString("aiPrompt.day")
}

func getAiAkConf() string {
	return viper.GetString("ai.ak")
}

func getAiModelConf() string {
	return viper.GetString("ai.model")
}

func getAiUrlConf() string {
	return viper.GetString("ai.url")
}

func getGitUsernameConf() string {
	return viper.GetString("git.username")
}

func getReportModeConf() string {
	return viper.GetString("report.mode")
}

func getGitReposConf() []string {
	return viper.GetStringSlice("git.repo")
}

func getReportIntervalDayConf() int {
	return viper.GetInt("report.intervalDay")
}

func getReportOutConf() string {
	return viper.GetString("report.out")
}
