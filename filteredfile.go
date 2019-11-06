package main

import (
	"fmt"
	"io"
	"strings"
)

type lineDetail struct {
	lineIndex uint
	posInFile int64
}

type filteredFile struct {
	Lines map[uint]lineDetail
}

func newFilteredFile(rs io.ReadSeeker, filter func(line FileLine) bool) *filteredFile {
	result := &filteredFile{}
	result.Lines = make(map[uint]lineDetail)
	tf := NewTextFile(rs, 1)

	for n, l := range tf.CachedLines {
		if filter(l) {
			result.Lines[n] = lineDetail{lineIndex: n, posInFile: l.position}
		}
	}
	return result
}

func (ff filteredFile) String() string {
	var sb strings.Builder
	for _, l := range ff.Lines {
		sb.WriteString(fmt.Sprintln("l:", l.lineIndex, "p:", l.posInFile))
	}
	return sb.String()
}
