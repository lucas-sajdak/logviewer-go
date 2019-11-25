package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

type FilteredLineTestCase struct {
	searchString string
	Lines        map[uint]int64
	firstLine    uint
	cacheSize    uint
}

func performFilteredFileTests(t *testing.T, testCases []FilteredLineTestCase, file io.ReadSeeker) {
	for n, c := range testCases {
		ff := newFilteredFile(file, c.cacheSize, func(line FileLine) bool {
			return strings.Contains(line.Contents, string(c.searchString))
		})
		ff.goTo(c.firstLine)
		if !reflect.DeepEqual(c.Lines, ff.Lines) {
			t.Errorf("Case %v: expect: %v have: %v", n, c.Lines, ff.Lines)
		}
	}
}

func TestFilteringSingleLineFile(t *testing.T) {
	f := newFileMock("base\n")
	testCases := []FilteredLineTestCase{
		{
			searchString: "bas",
			firstLine:    0,
			cacheSize:    1,
			Lines: map[uint]int64{
				0: 0}},
		{
			searchString: "base",
			firstLine:    0,
			cacheSize:    1,
			Lines: map[uint]int64{
				0: 0}},
		{
			searchString: "ba1e",
			firstLine:    0,
			cacheSize:    1,
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}

func TestFilteringUnfinishedLineFile(t *testing.T) {
	f := newFileMock("base")
	testCases := []FilteredLineTestCase{
		{
			searchString: "bas",
			firstLine:    0,
			cacheSize:    1,
			Lines:        map[uint]int64{}},
		{
			searchString: "base",
			firstLine:    0,
			cacheSize:    1,
			Lines:        map[uint]int64{}},
		{
			searchString: "ba1e",
			firstLine:    0,
			cacheSize:    1,
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}

func TestFilteringMultilineFile(t *testing.T) {
	f := newFileMock("text\nanotherThing\nanothertext\nsomethingelse\n")
	testCases := []FilteredLineTestCase{
		{
			searchString: "text",
			firstLine:    0,
			cacheSize:    1,
			Lines: map[uint]int64{
				0: 0}},
		{
			searchString: "text",
			firstLine:    0,
			cacheSize:    2,
			Lines: map[uint]int64{
				0: 0,
				2: 18}},
		{
			searchString: "e",
			firstLine:    0,
			cacheSize:    4,
			Lines: map[uint]int64{
				0: 0,
				1: 5,
				2: 18,
				3: 30}},
		{
			searchString: "e",
			firstLine:    2,
			cacheSize:    4,
			Lines: map[uint]int64{
				2: 18,
				3: 30}},
		{
			searchString: "Text",
			firstLine:    0,
			cacheSize:    1,
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}

func TestFilteringMultilineFileWithUnfinishedLastLine(t *testing.T) {
	f := newFileMock("complete\ncomplete\ncomplete")
	testCases := []FilteredLineTestCase{
		{
			searchString: "complete",
			firstLine:    0,
			cacheSize:    1,
			Lines: map[uint]int64{
				0: 0}},
		{
			searchString: "complete",
			firstLine:    1,
			cacheSize:    1,
			Lines: map[uint]int64{
				1: 9}},
		{
			searchString: "complete",
			firstLine:    0,
			cacheSize:    2,
			Lines: map[uint]int64{
				0: 0,
				1: 9}},
		{
			searchString: "complete",
			firstLine:    3,
			cacheSize:    1,
			Lines:        map[uint]int64{}},
		{
			searchString: "\n",
			firstLine:    0,
			cacheSize:    2,
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}
