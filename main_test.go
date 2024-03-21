package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

func TestGetEndGitUsername(t *testing.T) {
	_ = os.Setenv("GIT.USERNAME", "test")
	username := getEndGitUsername()
	expectedUsername := "test"
	if username != expectedUsername {
		t.Errorf("Expected %s, but got %s", expectedUsername, username)
	}
	assert.Equal(t, expectedUsername, username)
}

func TestGetGitLogAfter(t *testing.T) {
	repoPath := "C:\\Web\\zhuzhu-game\\.git"
	gitLog := "This is a git log message"
	gitLogMap := sync.Map{}
	maxLogLen := int64(0)

	getGitLogAfter(repoPath, &gitLog, &gitLogMap, &maxLogLen)

	actual, _ := gitLogMap.Load("zhuzhu-game")
	assert.Equal(t, gitLog, actual)
	assert.Equal(t, int64(len(gitLog)+len("zhuzhu-game")), maxLogLen)
}

func TestMakeAiReq(t *testing.T) {
	// 创建一个 sync.Map 并添加测试数据
	gitLogMap := &sync.Map{}
	gitLogMap.Store("repo1", "git log 1")
	gitLogMap.Store("repo2", "git log 2")

	// 设置 maxLogLen 和 prompt
	maxLogLen := int64(100)
	prompt := "prompt"

	// 调用被测试函数
	result := makeAiReq(gitLogMap, maxLogLen, prompt)

	// 验证结果是否符合预期
	expected := "\n以下是repo1仓库的GIT日志\ngit log 1\n以下是repo2仓库的GIT日志\ngit log 2"

	assert.Equal(t, expected, result)
}
