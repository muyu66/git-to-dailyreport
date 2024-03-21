package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test1(t *testing.T) {
	expectedResult := 1

	result := 1

	assert.Equal(t, expectedResult, result)
}
