package main

import "io"

func newFileMock(contents string) *fileMock {
	result := &fileMock{}
	result.contents = string(contents)
	return result
}

func createEmptyFile() io.ReadSeeker {
	return newFileMock("")
}

func create10DigitsFile() io.ReadSeeker {
	return newFileMock("0123456789")
}

func createMultiLineFile() io.ReadSeeker {
	return newFileMock("1\n2\n3\n")
}
