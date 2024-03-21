package main

import "github.com/spf13/viper"

func getAiPromptConf() string {
	return viper.GetString("aiPrompt")
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
