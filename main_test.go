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

func TestGetNewPrUrl(t *testing.T) {
	assert.Equal(t,
		"https://github.com/kyoshidajp/gh-browse-pr/compare/test?expand=1",
		GetNewPrUrl("https://github.com/kyoshidajp/gh-browse-pr", "test"),
	)
}

func TestIsNumberString(t *testing.T) {
	assert.Equal(t, true, IsNumberString("100"))
	assert.Equal(t, true, IsNumberString("001"))
	assert.Equal(t, true, IsNumberString("-1")) // allow
	assert.Equal(t, false, IsNumberString("test"))
	assert.Equal(t, false, IsNumberString("1test"))
}
