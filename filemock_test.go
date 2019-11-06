package main

import (
	"bytes"
	"io"
	"testing"
)

type seeking struct {
	expectedErr bool
	offset      int64
	whence      int
	expectedPos int64
}

type reading struct {
	expectedErr     bool
	bufSize         int
	expectedBuf     []byte
	expectedBufSize int
}

func performTests(t *testing.T, file io.ReadSeeker, operations []interface{}) {
	for n, op := range operations {
		switch v := op.(type) {
		case seeking:
			doSeeking(t, n, v, file)
		case reading:
			doReading(t, n, v, file)
		}
	}
}

func doSeeking(t *testing.T, opIndex int, op seeking, s io.Seeker) {
	observedPos, err := s.Seek(op.offset, op.whence)
	if observedErr := err != nil; observedErr != op.expectedErr {
		t.Errorf("Operation %v:, Seek() error. want: %v, have: %v", opIndex, op.expectedErr, observedErr)
	}
	if op.expectedErr == false && observedPos != op.expectedPos {
		t.Errorf("Operation %v: Seek() returned unexpected value. want: %v, have: %v", opIndex, op.expectedPos, observedPos)
	}
}

func doReading(t *testing.T, opIndex int, op reading, r io.Reader) {
	b := make([]byte, op.bufSize)
	nBytes, err := r.Read(b)
	if observedErr := err != nil; observedErr != op.expectedErr {
		t.Errorf("operation %v: Read() error. want: %v, have: %v", opIndex, op.expectedErr, observedErr)
	}
	if expectedBufSize := len(op.expectedBuf); nBytes != len(op.expectedBuf) {
		t.Errorf("Operation %v: Read() returned unexpected value. want: %v, have: %v", opIndex, expectedBufSize, nBytes)
	}
	if !bytes.Equal(b[:nBytes], op.expectedBuf) {
		t.Errorf("Operation %v: Read() returned unexpected value. want: %v, have: %v", opIndex, op.expectedBuf, b[:nBytes])
	}
}

func TestSeekingEmptyFile(t *testing.T) {
	file := createEmptyFile()
	operations := []interface{}{
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},   // 0
		seeking{expectedErr: true, offset: -1},                                         // 0
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},   // 0
		seeking{expectedErr: true, offset: -10},                                        // 0
		seeking{expectedErr: false, offset: 5, whence: io.SeekStart, expectedPos: 5},   // 5
		seeking{expectedErr: true, offset: -1},                                         //
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},   //
		seeking{expectedErr: false, offset: 10, whence: io.SeekStart, expectedPos: 10}, // 10
		seeking{expectedErr: false, offset: 10, whence: io.SeekEnd, expectedPos: 10},   // 10
		seeking{expectedErr: true, offset: -15, whence: io.SeekCurrent},                // 0 - error
		seeking{expectedErr: false, offset: 0, whence: io.SeekEnd, expectedPos: 0},     // 3
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},   // 0
		seeking{expectedErr: true, offset: -1, whence: io.SeekCurrent},                 // 0 - error
	}

	performTests(t, file, operations)
}

func TestReadingFromEmptyFile(t *testing.T) {
	file := createEmptyFile()
	operations := []interface{}{
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},
		reading{expectedErr: true},
		seeking{expectedErr: false, offset: 10, whence: io.SeekCurrent, expectedPos: 10},
		reading{expectedErr: true},
	}
	performTests(t, file, operations)
}

func TestSeeking10DigitsFile(t *testing.T) {
	file := create10DigitsFile()
	operations := []interface{}{
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},   // 0
		seeking{expectedErr: false, offset: 2, whence: io.SeekCurrent, expectedPos: 2}, // 2
		seeking{expectedErr: false, offset: 2, whence: io.SeekCurrent, expectedPos: 4}, // 4
		seeking{expectedErr: true, offset: -5},                                         // 0 - error
	}

	performTests(t, file, operations)
}

func TestReadingFrom10DigitsFile(t *testing.T) {
	file := create10DigitsFile()
	operations := []interface{}{
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0}, // 0
		reading{expectedErr: false, bufSize: 3, expectedBuf: []byte("012")},
		reading{expectedErr: false, bufSize: 3, expectedBuf: []byte("345")},
		seeking{expectedErr: false, offset: 2, whence: io.SeekCurrent, expectedPos: 8},  // 2
		seeking{expectedErr: false, offset: 2, whence: io.SeekCurrent, expectedPos: 10}, // 4
		seeking{expectedErr: true, offset: -5},                                          // 0 - error
	}
	performTests(t, file, operations)
}

func TestReadingFromMultiLineFile(t *testing.T) {
	file := createMultiLineFile()
	operations := []interface{}{
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},
		reading{expectedErr: false, bufSize: 2, expectedBuf: []byte("1\n")},
		reading{expectedErr: false, bufSize: 2, expectedBuf: []byte("2\n")},
		reading{expectedErr: false, bufSize: 2, expectedBuf: []byte("3\n")},
		reading{expectedErr: true, bufSize: 2},
		seeking{expectedErr: false, offset: 2, whence: io.SeekStart, expectedPos: 2},
		reading{expectedErr: false, bufSize: 6, expectedBuf: []byte("2\n3\n")},
		seeking{expectedErr: false, offset: 0, whence: io.SeekStart, expectedPos: 0},
		reading{expectedErr: false, bufSize: 10, expectedBuf: []byte("1\n2\n3\n")},
		reading{expectedErr: true, bufSize: 10},
		seeking{expectedErr: false, offset: 0, whence: io.SeekCurrent, expectedPos: 6},
	}
	performTests(t, file, operations)
}
