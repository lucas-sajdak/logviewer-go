package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// FileLine holds line of text from the file and its position inside the file
type FileLine struct {
	Contents string
	position int64
}

// TextFile keeps line of file in user-defined size cache
type TextFile struct {
	CachedLines       map[uint]FileLine
	cacheSize         uint // number of lines to be cached
	rs                io.ReadSeeker
	startingLineIndex uint
}

// NewTextFile creates new text file for given filepath
func NewTextFile(rs io.ReadSeeker, cacheSize uint) *TextFile {
	result := &TextFile{}
	result.rs = rs
	result.startingLineIndex = 0
	result.CachedLines = make(map[uint]FileLine)
	result.cacheSize = cacheSize
	result.rs.Seek(0, io.SeekStart)
	result.goTo(result.startingLineIndex)
	return result
}

func getLinePosition(rs io.ReadSeeker, fromLine uint, fromPos int64, lineOffset int) (int64, error) {
	if fromPos == 0 {
		return 0, nil
	}
	bufferSize := Min(1024*512, fromPos-1)
	rs.Seek(fromPos-1, io.SeekStart)
	rs.Seek(-bufferSize, io.SeekCurrent)
	r := bufio.NewReader(rs)

	b := make([]byte, bufferSize)
	if _, err := r.Read(b); err != nil {
		fmt.Println("Read() failed: ", err.Error())
	}

	var i int64
	currentLine := fromLine
	newLineOffset := lineOffset
	var posInFile int64
	for i = bufferSize - 1; i != -1; i-- {
		posInFile = fromPos - (bufferSize - i)
		if b[i] == '\n' {
			currentLine--
			newLineOffset++
			//			fmt.Println("Line at:", currentLine, newLineOffset, posInFile)
			if int(currentLine) == int(fromLine)+lineOffset {
				return posInFile, nil
			}
		}
	}

	if posInFile == 0 {
		return 0, fmt.Errorf("Error")
	}

	return getLinePosition(rs, currentLine, posInFile, newLineOffset)
}

func (tf *TextFile) goTo(lineIndex uint) {
	var curLine uint
	var p int64
	r := bufio.NewReader(tf.rs)

	if lineIndex < tf.startingLineIndex {
		p, _ = getLinePosition(
			tf.rs,
			tf.startingLineIndex,
			tf.CachedLines[tf.startingLineIndex].position,
			int(lineIndex)-int(tf.startingLineIndex))
		//		fmt.Println("Found at:", p)
		curLine = lineIndex

	} else {
		if _, ok := tf.CachedLines[lineIndex]; ok {
			line, _ := tf.CachedLines[lineIndex]
			curLine = lineIndex
			p = line.position
		} else {
			if _, ok := tf.CachedLines[tf.startingLineIndex+tf.cacheSize-1]; ok {
				line, _ := tf.CachedLines[tf.startingLineIndex+tf.cacheSize-1]
				curLine = tf.startingLineIndex + tf.cacheSize - 1
				p = line.position
			}
		}
	}

	tf.CachedLines = make(map[uint]FileLine)
	tf.rs.Seek(p, io.SeekStart)
	for {
		b, err := r.ReadBytes('\n')
		if err != nil {
			break
		}

		if curLine >= lineIndex {
			notRNEndLine := strings.TrimSuffix(string(b), "\r\n") // deal with "\r\n"
			notREndLine := strings.TrimSuffix(notRNEndLine, "\n") // deal with "\n"
			tf.CachedLines[curLine] = FileLine{notREndLine, p}
		}

		p += int64(len(b))

		curLine++
		if curLine > lineIndex+tf.cacheSize-1 {
			break
		}
	}

	tf.startingLineIndex = lineIndex
}

func (tf TextFile) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%T", tf))
	switch v := tf.rs.(type) {
	case *os.File:
		sb.WriteString(fmt.Sprintf(" (filename:%v) ", v.Name()))
	default:
		sb.WriteString(fmt.Sprintf(" (filetype:%T) ", v))
	}
	sb.WriteString(fmt.Sprintf("startingLine:%v cacheSize:%v", tf.startingLineIndex, tf.cacheSize))
	for k, v := range tf.CachedLines {
		sb.WriteString(fmt.Sprintf("L%v:%v(p:%v)\n", k, v.Contents, v.position))
	}
	sb.WriteString(fmt.Sprintf("%T-END", tf))
	return sb.String()
}
