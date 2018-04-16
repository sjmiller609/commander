package utils

import (
	"errors"
	"strings"
	"testing"
)

func TestEnsurePrefix(t *testing.T) {
	test1 := "1000"
	result1 := EnsurePrefix(test1, ":")
	if !strings.HasPrefix(result1, ":") {
		t.Error(errors.New("test1 wasn't prefixed with \":\""))
	}

	test2 := ":1000"
	result2 := EnsurePrefix(test2, ":")
	if result2[:2] == "::" {
		t.Error(errors.New("test2 was double prefixed"))
	}
}