package main

import "testing"

func TestMaxAreOnLeft(t *testing.T) {

	maxesAreOnLeft := []struct {
		left     int64
		right    int64
		expected int64
	}{
		{1, 0, 1},
		{10, 0, 10},
		{11, 2, 11},
		{50, 40, 50},
		{50, 50, 50},
	}

	for _, s := range maxesAreOnLeft {
		if Max(s.left, s.right) != s.expected {
			t.Errorf("%v is not max of %v and %v", s.expected, s.left, s.right)
		}
	}
}

func TestMaxAreOnRight(t *testing.T) {

	maxesAreOnLeft := []struct {
		left     int64
		right    int64
		expected int64
	}{
		{1, 10, 10},
		{10, 20, 20},
		{11, 22, 22},
		{50, 400, 400},
		{400, 400, 400},
	}

	for _, s := range maxesAreOnLeft {
		if Max(s.left, s.right) != s.expected {
			t.Errorf("%v is not max of %v and %v", s.expected, s.left, s.right)
		}
	}
}

func TestMaxAreLeftAndRight(t *testing.T) {

	maxesAreOnLeft := []struct {
		left     int64
		right    int64
		expected int64
	}{
		{1, 1, 1},
		{10, 10, 10},
		{11, 11, 11},
		{50, 50, 50},
	}

	for _, s := range maxesAreOnLeft {
		if Max(s.left, s.right) != s.expected {
			t.Errorf("%v is not max of %v and %v", s.expected, s.left, s.right)
		}
	}
}
