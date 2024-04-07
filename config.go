package main

import (
	"github.com/spf13/viper"
	"strconv"
)

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

func getAiConf() AiConfig {
	usedAi := getUseAiConf()
	m := viper.GetStringMapString("ai." + usedAi)
	maxInputTokens, _ := strconv.ParseUint(m["maxinputtokens"], 10, 32)
	return AiConfig{
		Model:          m["model"],
		BaseUrl:        m["baseurl"],
		ApiKey:         m["apikey"],
		MaxInputTokens: uint32(maxInputTokens),
	}
}

func getUseAiConf() string {
	return viper.GetString("useAi")
}
