package main

import (
	"fmt"
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
			got := findUnique(test.vec)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestCheckSchedule(t *testing.T) {
	tests := map[string]struct {
		min     int
		max     int
		p       string
		partLim [2]int
		want    error
	}{
		"correct (minute)":     {min: 0, max: 5, p: "0-5", partLim: [2]int{0, 59}, want: nil},
		"min > max (weekDay)":  {min: 6, max: 5, p: "6-5", partLim: [2]int{0, 6}, want: fmt.Errorf("The ranges must be defined as 'min-max' with `min` <= `max`. Expects %v <= %v from string %s", 6, 5, "6-5")},
		"min lower (monthDay)": {min: 0, max: 25, p: "0-25", partLim: [2]int{1, 31}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v", 1, 31, 0, 25)},
		"max higher (month)":   {min: 5, max: 13, p: "5-13", partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v", 1, 12, 5, 13)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := CheckSchedule(test.min, test.max, test.p, test.partLim)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

//func ParseSchedule(schedule string) (Schedule, error)
func TestParseSchedule(t *testing.T) {

	every := ParsedPart{
		Any: true,
	}

	everySecond := Schedule{
		Second:   every,
		Minute:   every,
		Hour:     every,
		WeekDay:  every,
		MonthDay: every,
		Month:    every,
	}

	minutesMonday := Schedule{
		Second:   ParsedPart{List: []int{0}},
		Minute:   ParsedPart{Min: 2, Max: 5},
		Hour:     every,
		WeekDay:  ParsedPart{List: []int{0}},
		MonthDay: every,
		Month:    every,
	}

	specific := Schedule{
		Second:   ParsedPart{List: []int{0}},
		Minute:   ParsedPart{List: []int{30}},
		Hour:     ParsedPart{List: []int{12}},
		WeekDay:  every,
		MonthDay: ParsedPart{List: []int{5}},
		Month:    ParsedPart{List: []int{1}},
	}

	listHours := Schedule{
		Second:   ParsedPart{List: []int{0}},
		Minute:   ParsedPart{List: []int{0}},
		Hour:     ParsedPart{List: []int{3, 5, 6}},
		WeekDay:  every,
		MonthDay: ParsedPart{List: []int{31}},
		Month:    every,
	}

	tests := map[string]struct {
		sch  string
		want Schedule
	}{
		"every second":                   {sch: "* * * * * *", want: everySecond},
		"range minutes on Monday":        {sch: "0 2-5 * 0 * *", want: minutesMonday},
		"specific time on 5th January ":  {sch: "0 30 12 * 5 1", want: specific},
		"list hours every 31th monthDay": {sch: "0 0 3,5,6 * 31 *", want: listHours},
		//TODO: many more tests!
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := ParseSchedule(test.sch)
			if !reflect.DeepEqual(test.want, got) {
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
