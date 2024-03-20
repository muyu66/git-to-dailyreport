package main

import "github.com/spf13/viper"

func getAiPrompt() string {
	return viper.GetString("aiPrompt")
}

func getAiAk() string {
	return viper.GetString("ai.ak")
}

func getAiModel() string {
	return viper.GetString("ai.model")
}

func getAiUrl() string {
	return viper.GetString("ai.url")
}

func getGitCommit() string {
	return viper.GetString("git.commit")
}

func getGitPath() string {
	return viper.GetString("git.path")
}

func getGitUsername() string {
	return viper.GetString("git.username")
}

func getReportMode() string {
	return viper.GetString("report.mode")
}
