package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
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
	file              *os.File
	startingLineIndex uint
}

// NewTextFile creates new text file for given filepath
func NewTextFile(filepath string, cacheSize uint) *TextFile {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModeCharDevice)
	if err != nil {
		log.Panicf("OpenFile() failed: %v", err.Error())
	}

	result := &TextFile{}
	result.file = f
	result.startingLineIndex = 0
	result.CachedLines = make(map[uint]FileLine)
	result.cacheSize = cacheSize
	result.file.Seek(0, io.SeekStart)
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

	//	fmt.Printf("read:<%s>", string(b))

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
	r := bufio.NewReader(tf.file)
	if lineIndex < tf.startingLineIndex {
		p, _ = getLinePosition(
			tf.file,
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
	tf.file.Seek(p, io.SeekStart)
	for {
		b, err := r.ReadBytes('\n')
		if err != nil {
			log.Panic("ReadBytes() failed: ", err.Error())
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
	sb.WriteString(fmt.Sprintf("Printing TextFile - file:%v startingLine:%v cacheSize:%v\n", tf.file.Name(), tf.startingLineIndex, tf.cacheSize))

	for k, v := range tf.CachedLines {
		sb.WriteString(fmt.Sprintf("L%v:%v(p:%v)\n", k, v.Contents, v.position))
	}
	sb.WriteString("Printing TextFile - END\n")
	return sb.String()
}
