package main

import (
	"reflect"
	"testing"
)

func TestFindMin(t *testing.T) {
	tests := map[string]struct {
		vec  []int
		want int
	}{
		"0 to 5":  {vec: []int{0, 1, 2, 3, 4, 5}, want: 0},
		"-2 to 2": {vec: []int{-2, -1, 0, 1, 2}, want: -2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := findMin(test.vec)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestFindMax(t *testing.T) {
	tests := map[string]struct {
		vec  []int
		want int
	}{
		"0 to 5":  {vec: []int{0, 1, 2, 3, 4, 5}, want: 5},
		"-2 to 2": {vec: []int{-2, -1, 0, 1, 2}, want: 2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := findMax(test.vec)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestMakeRange(t *testing.T) {
	tests := map[string]struct {
		min  int
		max  int
		want []int
		err  bool
	}{
		"0 to 5": {min: 0, max: 5, want: []int{0, 1, 2, 3, 4, 5}, err: false},
		"9 to 5": {min: 9, max: 5, want: []int{}, err: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := makeRange(test.min, test.max)
			if !reflect.DeepEqual(test.want, got) && !reflect.DeepEqual(test.err, err != nil) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestCompareTime(t *testing.T) {
	tests := map[string]struct {
		sched     ParsedPart
		dt        int
		want      int
		wantShift int
	}{
		"Any":     {sched: ParsedPart{Any: true}, dt: 0, want: 0, wantShift: 0},
		"List":    {sched: ParsedPart{List: []int{0, 2, 4}}, dt: 1, want: 2, wantShift: 0},
		"Min-Max": {sched: ParsedPart{Min: 4, Max: 9}, dt: 10, want: 4, wantShift: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, shift := CompareTime(test.sched, test.dt)
			if !reflect.DeepEqual(test.want, got) && !reflect.DeepEqual(test.wantShift, shift) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}
