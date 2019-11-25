package main

import (
	"fmt"
	"io"
	"strings"
)

type filteredFile struct {
	Lines          map[ /*lineIndex*/ uint] /*position in file*/ int64
	filter         func(FileLine) bool
	firstLineIndex uint
	cacheSize      uint
	tf             *TextFile
}

func newFilteredFile(rs io.ReadSeeker, cacheSize uint, filter func(line FileLine) bool) *filteredFile {
	result := &filteredFile{}
	result.firstLineIndex = 0
	result.cacheSize = cacheSize
	result.filter = filter
	result.tf = NewTextFile(rs, 1)
	result.goTo(0)
	return result
}

func (ff *filteredFile) goTo(firstLine uint) {
	ff.Lines = make(map[uint]int64)
	ff.tf.goTo(firstLine)

	keepFiltering := true
	checkingLine := firstLine
	for keepFiltering {
		var i uint = 0
		for ; i != ff.tf.cacheSize; i++ {
			if _, ok := ff.tf.CachedLines[checkingLine]; ok {
				line, _ := ff.tf.CachedLines[checkingLine]
				if ff.filter(line) {
					ff.Lines[checkingLine] = line.position
				}
			}

			checkingLine++

			if uint(len(ff.Lines)) >= ff.cacheSize {
				keepFiltering = false
				break
			}
		}

		ff.tf.goTo(checkingLine)
		if len(ff.tf.CachedLines) == 0 {
			keepFiltering = false
		}
	}
}

func (ff filteredFile) String() string {
	var sb strings.Builder
	for n, l := range ff.Lines {
		sb.WriteString(fmt.Sprintln("l:", n, "p:", l))
	}
	return sb.String()
}
