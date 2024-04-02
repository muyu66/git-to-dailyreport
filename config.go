package main

import "github.com/spf13/viper"

func getAiNameConf() string {
	return viper.GetString("ai.name")
}

func getAiAkConf() string {
	return viper.GetString("ai.ak")
}

func getAiModelConf() string {
	return viper.GetString("ai.model")
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

func getReportLangConf() string {
	return viper.GetString("report.lang")
}

func getReportFlowConf() bool {
	return viper.GetBool("report.flow")
}
