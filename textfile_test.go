package main

import (
	"io"
	"reflect"
	"testing"
)

type textFileTestCase struct {
	startingLineIndex uint
	CacheSize         uint
	Lines             map[uint]FileLine
}

func performTextFileTests(t *testing.T, testCases []textFileTestCase, rs io.ReadSeeker) {
	for n, c := range testCases {
		tf := NewTextFile(rs, c.CacheSize)
		tf.goTo(c.startingLineIndex)
		if !reflect.DeepEqual(tf.CachedLines, c.Lines) {
			t.Errorf("Op %v: expect: %v, have: %v", n, tf.CachedLines, c.Lines)
		}
	}
}

func TestMultilineTextFile(t *testing.T) {
	rs := newFileMock("1st\n2nd\n3rd\n")

	testCases := []textFileTestCase{
		{
			startingLineIndex: 0,
			CacheSize:         1,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0}}},
		{
			startingLineIndex: 0,
			CacheSize:         3,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0},
				1: FileLine{Contents: "2nd", position: 4},
				2: FileLine{Contents: "3rd", position: 8}}},
		{
			startingLineIndex: 0,
			CacheSize:         5,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0},
				1: FileLine{Contents: "2nd", position: 4},
				2: FileLine{Contents: "3rd", position: 8}}},
		{
			startingLineIndex: 0,
			CacheSize:         2,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0},
				1: FileLine{Contents: "2nd", position: 4}}},
		{
			startingLineIndex: 1,
			CacheSize:         2,
			Lines: map[uint]FileLine{
				1: FileLine{Contents: "2nd", position: 4},
				2: FileLine{Contents: "3rd", position: 8}}},
		{
			startingLineIndex: 4,
			CacheSize:         2,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 1,
			CacheSize:         1,
			Lines: map[uint]FileLine{
				1: FileLine{Contents: "2nd", position: 4}}},
	}

	performTextFileTests(t, testCases, rs)
}

func TestSingleLineTextFile(t *testing.T) {
	rs := newFileMock("1st\n")

	testCases := []textFileTestCase{
		{
			startingLineIndex: 0,
			CacheSize:         1,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0}}},
		{
			startingLineIndex: 0,
			CacheSize:         2,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0}}},
		{
			startingLineIndex: 0,
			CacheSize:         3,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0}}},
		{
			startingLineIndex: 1,
			CacheSize:         1,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 1,
			CacheSize:         2,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 1,
			CacheSize:         3,
			Lines:             map[uint]FileLine{}},
	}

	performTextFileTests(t, testCases, rs)
}

func TestUnfinishedLineTextFile(t *testing.T) {
	rs := newFileMock("1st")

	testCases := []textFileTestCase{
		{
			startingLineIndex: 0,
			CacheSize:         1,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 0,
			CacheSize:         2,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 0,
			CacheSize:         3,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 1,
			CacheSize:         1,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 1,
			CacheSize:         2,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 1,
			CacheSize:         3,
			Lines:             map[uint]FileLine{}},
	}

	performTextFileTests(t, testCases, rs)
}

func TestUnfinishedMultilineTextFile(t *testing.T) {
	rs := newFileMock("1st\n2nd\n3rd")

	testCases := []textFileTestCase{
		{
			startingLineIndex: 0,
			CacheSize:         1,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0}}},
		{
			startingLineIndex: 0,
			CacheSize:         2,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0},
				1: FileLine{Contents: "2nd", position: 4}}},
		{
			startingLineIndex: 0,
			CacheSize:         3,
			Lines: map[uint]FileLine{
				0: FileLine{Contents: "1st", position: 0},
				1: FileLine{Contents: "2nd", position: 4}}},
		{
			startingLineIndex: 1,
			CacheSize:         1,
			Lines: map[uint]FileLine{
				1: FileLine{Contents: "2nd", position: 4}}},
		{
			startingLineIndex: 1,
			CacheSize:         2,
			Lines: map[uint]FileLine{
				1: FileLine{Contents: "2nd", position: 4}}},

		{
			startingLineIndex: 2,
			CacheSize:         1,
			Lines:             map[uint]FileLine{}},
		{
			startingLineIndex: 2,
			CacheSize:         2,
			Lines:             map[uint]FileLine{}},
	}

	performTextFileTests(t, testCases, rs)
}
