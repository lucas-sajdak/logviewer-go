package main

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// fileMock is a mockup for io.File struct
type fileMock struct {
	io.ReadSeeker

	contents string
	position int64
}

type seeking struct {
	offset      int64
	whence      int
	expectedErr bool
	expectedPos int64
}

func newFile(contents []byte) *fileMock {
	result := &fileMock{}
	result.contents = string(contents)
	return result
}

func (fm fileMock) Seek(offset int64, whence int) (pos int64, err error) {
	err = nil
	switch whence {
	case io.SeekStart:
		if offset < 0 {
			err = errors.New("Position before start of the file")
		}
		fm.position = offset
		pos = offset
	}
	return pos, err
}

func (fm fileMock) Read(p []byte) (n int, err error) {

	if int(fm.position) >= len(fm.contents) {
		err = io.EOF
		return
	}

	err = nil
	availableSize := len(fm.contents) - int(fm.position)
	copyingSize := int(Min(int64(availableSize), int64(len(p))))
	copy(p[0:copyingSize], []byte(fm.contents)[fm.position:fm.position+int64(copyingSize)])

	//	if int(fm.position)+copyingSize == len(fm.contents) {
	//		err = io.EOF
	//	}

	n = copyingSize
	return n, err
}

func CreateEmptyFile() io.ReadSeeker {
	return newFile([]byte(""))
}

func Create10DigitsFile() io.ReadSeeker {
	return newFile([]byte("0123456789"))
}

func TestSeekingEmptyFile(t *testing.T) {
	file := CreateEmptyFile()

	seekings := []seeking{
		{offset: 0, whence: io.SeekStart, expectedErr: false, expectedPos: 0},
		{offset: -1, whence: io.SeekStart, expectedErr: true},
		{offset: 0, whence: io.SeekStart, expectedErr: false, expectedPos: 0},
		{offset: -10, whence: io.SeekStart, expectedErr: true},
		{offset: 0, whence: io.SeekStart, expectedErr: false, expectedPos: 0},
		{offset: 1, whence: io.SeekStart, expectedErr: false, expectedPos: 1},
		{offset: 10, whence: io.SeekStart, expectedErr: false, expectedPos: 10},
		{offset: 0, whence: io.SeekStart, expectedErr: false, expectedPos: 0},
	}

	for n, s := range seekings {
		pos, err := file.Seek(s.offset, s.whence)
		if observed := err != nil; observed != s.expectedErr {
			t.Errorf("%v-Observed: %v (%v), expected: %v", n, observed, err, s.expectedErr)
		}
		if err == nil && pos != s.expectedPos {
			t.Errorf("%v-Observed: %v, expected: %v", n, pos, s.expectedPos)
		}
	}
}

func TestReadingFromEmptyFile(t *testing.T) {
	file := CreateEmptyFile()
	readings := []struct {
		readingPos     int64
		size           int
		expectedErr    bool
		expectedBuffer []byte
	}{
		{size: 1, expectedErr: true},
		{size: 5, expectedErr: true},
		{size: 10, expectedErr: true},
		{size: 20, expectedErr: true},
	}

	for _, r := range readings {
		b := make([]byte, r.size)
		_, err := file.Read(b)
		if err != nil {
			if observedErr := err != nil; observedErr != r.expectedErr {
				t.Errorf("read err")
			}
			continue
		}

		if !bytes.Equal(b, r.expectedBuffer) {
			t.Errorf("Read unexpected contents, want: '%s', have: '%s'", b, r.expectedBuffer)
		}
	}
}

func TestSeeking10DigitsFile(t *testing.T) {
	file := Create10DigitsFile()
	seekings := []seeking{
		{offset: 0, whence: io.SeekStart, expectedErr: false, expectedPos: 0},
		{offset: 1, whence: io.SeekStart, expectedErr: false, expectedPos: 1},
		{offset: 2, whence: io.SeekStart, expectedErr: false, expectedPos: 2},
		{offset: 3, whence: io.SeekStart, expectedErr: false, expectedPos: 3},
		{offset: 9, whence: io.SeekStart, expectedErr: false, expectedPos: 9},
		{offset: -1, whence: io.SeekStart, expectedErr: true},
		{offset: 0, whence: io.SeekStart, expectedErr: false, expectedPos: 0},
		{offset: -10, whence: io.SeekStart, expectedErr: true},
		{offset: 1, whence: io.SeekStart, expectedErr: false, expectedPos: 1},
		{offset: 2, whence: io.SeekStart, expectedErr: false, expectedPos: 2},
		{offset: 10, whence: io.SeekStart, expectedErr: false, expectedPos: 10},
		{offset: 7, whence: io.SeekStart, expectedErr: false, expectedPos: 7},
	}

	for n, s := range seekings {
		pos, err := file.Seek(s.offset, s.whence)

		if observed := err != nil; observed != s.expectedErr {
			t.Errorf("%v-Observed: %v (%v), expected: %v", n, observed, err, s.expectedErr)
		}

		if err == nil && pos != s.expectedPos {
			t.Errorf("%v-Observed: %v, expected: %v", n, pos, s.expectedPos)
		}
	}
}

func TestReadingFrom10DigitsFile(t *testing.T) {
	file := Create10DigitsFile()

	readings := []struct {
		readingPos     int64
		size           int
		expectedErr    bool
		expectedBuffer []byte
	}{
		{size: 1, expectedBuffer: []byte("0")},
		{size: 5, expectedBuffer: []byte("01234")},
		{size: 10, expectedBuffer: []byte("0123456789")},
		{size: 20, expectedBuffer: []byte("0123456789")},
	}

	for n, r := range readings {
		b := make([]byte, r.size)
		nBytes, err := file.Read(b)
		if err != nil {
			t.Errorf("%v, read err", n)
		}

		if !bytes.Equal(b[:nBytes], r.expectedBuffer) {
			t.Errorf("%v, Read unexpected contents, want: %v, have: %v", n, b, r.expectedBuffer)
		}
	}
}
