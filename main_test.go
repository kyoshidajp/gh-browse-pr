package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrUrl(t *testing.T) {
	expected := "https://github.com/kyoshidajp/gh-browse-pr/pull/main"
	actual := GetPrUrl(
		"https://github.com/kyoshidajp/gh-browse-pr",
		"main",
	)

	assert.Equal(t, expected, actual)
}
