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
}

func performFilteredFileTests(t *testing.T, testCases []FilteredLineTestCase, file io.ReadSeeker) {
	for n, c := range testCases {
		ff := newFilteredFile(file, func(line FileLine) bool {
			return strings.Contains(line.Contents, string(c.searchString))
		})

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
			Lines: map[uint]int64{
				0: 0}},
		{
			searchString: "base",
			Lines: map[uint]int64{
				0: 0}},
		{
			searchString: "ba1e",
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}

func TestFilteringUnfinishedLineFile(t *testing.T) {
	f := newFileMock("base")
	testCases := []FilteredLineTestCase{
		{
			searchString: "bas",
			Lines:        map[uint]int64{}},
		{
			searchString: "base",
			Lines:        map[uint]int64{}},
		{
			searchString: "ba1e",
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}

func TestFilteringMultilineFile(t *testing.T) {
	f := newFileMock("text\nanotherThing\nanothertext\nsomethingelse\n")
	testCases := []FilteredLineTestCase{
		{
			searchString: "text",
			Lines: map[uint]int64{
				0: 0,
				2: 18}},
		{
			searchString: "e",
			Lines: map[uint]int64{
				0: 0,
				1: 5,
				2: 18,
				3: 30}},
		{
			searchString: "Text",
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}

func TestFilteringMultilineFileWithUnfinishedLastLine(t *testing.T) {
	f := newFileMock("complete\ncomplete\ncomplete")
	testCases := []FilteredLineTestCase{
		{
			searchString: "complete",
			Lines: map[uint]int64{
				0: 0,
				1: 9}},
		{
			searchString: "\n",
			Lines:        map[uint]int64{}},
	}
	performFilteredFileTests(t, testCases, f)
}
