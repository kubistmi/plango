package utils

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
			got := FindMin(test.vec)
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
			got := FindMax(test.vec)
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
			got, err := MakeRange(test.min, test.max)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
			if !reflect.DeepEqual(test.err, err != nil) {
				t.Fatalf("Expected no error, got: %#v", err)
			}
		})
	}
}

func TestFindUnique(t *testing.T) {
	tests := map[string]struct {
		vec  []int
		want []int
	}{
		"all unique":     {vec: []int{0, 1, 2, 3, 4}, want: []int{0, 1, 2, 3, 4}},
		"all same":       {vec: []int{1, 1, 1, 1}, want: []int{1}},
		"few duplicates": {vec: []int{0, 0, 1, 4, 4, 4}, want: []int{0, 1, 4}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := FindUnique(test.vec)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestIsIn(t *testing.T) {
	tests := map[string]struct {
		element int
		vec     []int
		want    bool
	}{
		"its first":     {element: 6, vec: []int{6, 9, 7, 6, 5}, want: true},
		"its last":      {element: 1, vec: []int{7, 4, 3, 9, 7, 0, 1}, want: true},
		"its middle":    {element: -1, vec: []int{6, 4, 7, 9, -1, 5, 2, 4}, want: true},
		"its not there": {element: 99, vec: []int{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3}, want: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := IsIn(test.element, test.vec)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}
