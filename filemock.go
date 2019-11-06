package main

import (
	"errors"
	"io"
)

// fileMock is a mockup for io.File struct
type fileMock struct {
	io.ReadSeeker
	contents string
	position int64
}

func (fm *fileMock) Seek(offset int64, whence int) (int64, error) {
	var err error
	var globalPos int64
	switch whence {
	case io.SeekStart:
		globalPos = offset
	case io.SeekEnd:
		globalPos = int64(len(fm.contents)) + offset
	case io.SeekCurrent:
		globalPos = fm.position + offset
	}
	if globalPos < 0 {
		err = errors.New("Position before start of the file")
		fm.position = 0
	} else {
		err = nil
		fm.position = globalPos
	}
	return fm.position, err
}

func (fm *fileMock) Read(p []byte) (n int, err error) {
	if int(fm.position) >= len(fm.contents) {
		err = io.EOF
		return
	}

	err = nil
	availableSize := len(fm.contents) - int(fm.position)
	copyingSize := int(Min(int64(availableSize), int64(len(p))))
	copy(p[0:copyingSize], []byte(fm.contents)[fm.position:fm.position+int64(copyingSize)])
	fm.position += int64(copyingSize)
	n = copyingSize
	return n, err
}
