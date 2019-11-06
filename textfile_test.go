package main

import (
	"reflect"
	"testing"
)

type textFileTest struct {
	FileContent       string
	startingLineIndex uint
	cacheSize         uint
	Lines             map[uint]FileLine
}

func TestTextFile(t *testing.T) {

	test := textFileTest{
		FileContent: "test", cacheSize: 1}
	fm := newFileMock(test.FileContent)
	tf := NewTextFile(fm, test.cacheSize)
	tf.goTo(tf.startingLineIndex)

	if reflect.DeepEqual(tf.CachedLines, test.Lines) {
		t.Errorf("Op %v:", 1)
	}

}
