package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestFilteredFile(t *testing.T) {
	f := newFileMock("base\nbase\nbass\n")
	ff := newFilteredFile(f, func(line FileLine) bool {
		return strings.Contains(line.Contents, string("bas"))
	})
	fmt.Printf("%v", ff)
	t.Log(fmt.Sprintf("ff: %v", ff))
}
