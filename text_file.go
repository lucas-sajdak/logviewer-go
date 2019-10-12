package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// TextFile keeps line of file in user-defined size cache
type TextFile struct {
	Lines             map[int]string // already read lines
	cacheSize         int            // number of lines to be cached
	file              *os.File
	startingLineIndex int
}

// NewTextFile creates new text file for given filepath
func NewTextFile(filepath string, cachedLines int) *TextFile {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModeCharDevice)
	if err != nil {
		log.Panicf("OpenFile() failed: %v", err.Error())
	}

	result := &TextFile{}
	result.file = f
	result.startingLineIndex = 0
	result.cacheSize = cachedLines
	return result
}

func (tf *TextFile) gotoLine(index int) {
	tf.file.Seek(0, io.SeekStart)

	s := bufio.NewScanner(tf.file)
	curLine := 0
	tf.Lines = make(map[int]string)

	for s.Scan() {
		if curLine >= index && curLine < index+tf.cacheSize {
			tf.Lines[curLine] = s.Text()
		}

		if curLine > index+tf.cacheSize {
			break
		}

		curLine++
	}
	tf.startingLineIndex = index
}

func (tf TextFile) String() string {
	var sb strings.Builder
	sb.WriteString("Printing TextFile - BEGIN\n")
	sb.WriteString(fmt.Sprintf("Cached - file:%v startingLine:%v cacheSize:%v\n", tf.file.Name(), tf.startingLineIndex, tf.cacheSize))

	//	for i := tf.startingLineIndex; i < tf.startingLineIndex+tf.cacheSize; i++ {
	//		sb.WriteString(fmt.Sprintf("%v>>%v\n", i, tf.Lines[i]))
	//	}
	//	sb.WriteString("map:\n")

	for k, v := range tf.Lines {
		sb.WriteString(fmt.Sprintf("%v>>%v\n", k, v))
	}
	sb.WriteString("Printing TextFile - END\n")
	return sb.String()
}
