package main

import (
	"fmt"
	"io"
	"strings"
)

type filteredFile struct {
	Lines map[ /*lineIndex*/ uint] /*position in file*/ int64
}

func newFilteredFile(rs io.ReadSeeker, filter func(line FileLine) bool) *filteredFile {
	result := &filteredFile{}
	result.Lines = make(map[uint]int64)
	tf := NewTextFile(rs, 1)

	success := true
	for success {
		for n, l := range tf.CachedLines {
			if filter(l) {
				result.Lines[n] = l.position
			}
		}
		tf.goTo(tf.startingLineIndex + tf.cacheSize)
		success = len(tf.CachedLines) > 0
	}

	return result
}

func (ff filteredFile) String() string {
	var sb strings.Builder
	for n, l := range ff.Lines {
		sb.WriteString(fmt.Sprintln("l:", n, "p:", l))
	}
	return sb.String()
}
