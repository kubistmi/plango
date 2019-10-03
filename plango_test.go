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
			got := findUnique(test.vec)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestAnyCheckSchedule(t *testing.T) {

	partAll := PartAny{Text: "*"}

	tests := map[string]struct {
		part    PartAny
		partLim [2]int
		want    error
	}{
		"no error (ever)": {part: partAll, partLim: [2]int{0, 59}, want: nil},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.part.checkPart(test.partLim)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestIntervalCheckSchedule(t *testing.T) {

	int05 := PartInterval{Text: "0-5", Min: 0, Max: 5}
	int65 := PartInterval{Text: "6-5", Min: 6, Max: 5}
	int25 := PartInterval{Text: "0-25", Min: 0, Max: 25}
	int513 := PartInterval{Text: "5-13", Min: 5, Max: 13}

	tests := map[string]struct {
		part    SchedulePart
		partLim [2]int
		want    error
	}{
		"correct (minute)":     {part: int05, partLim: [2]int{0, 59}, want: nil},
		"min > max (weekDay)":  {part: int65, partLim: [2]int{0, 6}, want: fmt.Errorf("The ranges must be defined as 'min-max' with `min` <= `max`. Expects %v <= %v from string %s", 6, 5, int65.Text)},
		"min lower (monthDay)": {part: int25, partLim: [2]int{1, 31}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 31, 0, 25, int25.Text)},
		"max higher (month)":   {part: int513, partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 12, 5, 13, int513.Text)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.part.checkPart(test.partLim)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestListCheckSchedule(t *testing.T) {

	list50 := PartList{Text: "0,50", List: []int{0, 50}}
	list42 := PartList{Text: "4,2", List: []int{4, 2}}
	listSingle := PartList{Text: "23", List: []int{23}}
	list09 := PartList{Text: "0,9", List: []int{0, 9}}
	list513 := PartList{Text: "5,6,10,15", List: []int{5, 6, 10, 15}}

	tests := map[string]struct {
		part    SchedulePart
		partLim [2]int
		want    error
	}{
		"correct (minute)":    {part: list50, partLim: [2]int{0, 59}, want: nil},
		"min > max (weekDay)": {part: list42, partLim: [2]int{0, 6}, want: nil},
		"single value (hour)": {part: listSingle, partLim: [2]int{0, 23}, want: nil},
		"min lower (month)":   {part: list09, partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 12, 0, 9, list09.Text)},
		"max higher (month)":  {part: list513, partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 12, 5, 15, list513.Text)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.part.checkPart(test.partLim)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestParseSchedule(t *testing.T) {

	every := PartAny{Text: "*"}

	everySecond := Schedule{
		Second:   every,
		Minute:   every,
		Hour:     every,
		WeekDay:  every,
		MonthDay: every,
		Month:    every,
	}

	minutesMonday := Schedule{
		Second:   PartList{Text: "0", List: []int{0}},
		Minute:   PartInterval{Text: "2-5", Min: 2, Max: 5},
		Hour:     every,
		WeekDay:  PartList{Text: "0", List: []int{0}},
		MonthDay: every,
		Month:    every,
	}

	specific := Schedule{
		Second:   PartList{Text: "0", List: []int{0}},
		Minute:   PartList{Text: "30", List: []int{30}},
		Hour:     PartList{Text: "12", List: []int{12}},
		WeekDay:  every,
		MonthDay: PartList{Text: "5", List: []int{5}},
		Month:    PartList{Text: "1", List: []int{1}},
	}

	listHours := Schedule{
		Second:   PartList{Text: "0", List: []int{0}},
		Minute:   PartList{Text: "0", List: []int{0}},
		Hour:     PartList{Text: "3,5,6", List: []int{3, 5, 6}},
		WeekDay:  every,
		MonthDay: PartList{Text: "31", List: []int{31}},
		Month:    every,
	}

	tests := map[string]struct {
		sch  string
		want Schedule
		err  error
	}{
		"every second":                   {sch: "* * * * * *", want: everySecond, err: nil},
		"range minutes on Monday":        {sch: "0 2-5 * 0 * *", want: minutesMonday, err: nil},
		"specific time on 5th January ":  {sch: "0 30 12 * 5 1", want: specific, err: nil},
		"list hours every 31th monthDay": {sch: "0 0 3,5,6 * 31 *", want: listHours, err: nil},
		"error monthDay too high":        {sch: "0 0 12 * 32 *", want: Schedule{}, err: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 31, 32, 32, "32")},
		"error too many fields":          {sch: "0 0 0 0 0 0 *", want: Schedule{}, err: fmt.Errorf("Incorrect number of fields, expected 6 got %v. Fields are separated by a space and the whitespace can't be used for any other purpose", 7)},
		"error wrong range":              {sch: "0 0 12-18-10 0 0 *", want: Schedule{}, err: fmt.Errorf("Incorrect format of range. Expected 2 values separated by `-`, got %v", 3)},
		"error single non-convertibl":    {sch: "a b c d e *", want: Schedule{}, err: fmt.Errorf("Unable to parse part of schedule: %s", "a")},
		"error non-convertible range":    {sch: "0 0 0 0 12-18a *", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "18a")},
		"error non-convertible list":     {sch: "0 0 0 0 12,1a *", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "1a")},
		"error list and range":           {sch: "0 11-15,16 0 0 0 *", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "15,16")},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseSchedule(test.sch)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
			if !reflect.DeepEqual(test.err, err) {
				t.Fatalf("Expected: %#v, got: %#v", test.err, err)
			}
		})
	}
}

func TestAnyCompareTime(t *testing.T) {
	tests := map[string]struct {
		sched     SchedulePart
		dt        int
		want      int
		wantShift int
	}{
		"no-shift zero":       {sched: PartAny{Text: "*"}, dt: 0, want: 0, wantShift: 0},
		"no-shift fifty-nine": {sched: PartAny{Text: "*"}, dt: 59, want: 59, wantShift: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, shift := test.sched.compareTime(test.dt)
			if !reflect.DeepEqual(test.want, got) && !reflect.DeepEqual(test.wantShift, shift) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestIntervalCompareTime(t *testing.T) {
	tests := map[string]struct {
		sched     SchedulePart
		dt        int
		want      int
		wantShift int
	}{
		"no-shift in set": {sched: PartInterval{Text: "0-3", Min: 0, Max: 3}, dt: 2, want: 2, wantShift: 0},
		"no-shift lower":  {sched: PartInterval{Text: "0-3", Min: 5, Max: 10}, dt: 2, want: 5, wantShift: 0},
		"shift higher":    {sched: PartInterval{Text: "4-9", Min: 4, Max: 9}, dt: 10, want: 4, wantShift: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, shift := test.sched.compareTime(test.dt)
			if !reflect.DeepEqual(test.want, got) && !reflect.DeepEqual(test.wantShift, shift) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}

func TestListCompareTime(t *testing.T) {
	tests := map[string]struct {
		sched     PartList
		dt        int
		want      int
		wantShift int
	}{
		"no-shift exact": {sched: PartList{Text: "27", List: []int{27}}, dt: 27, want: 27, wantShift: 0},
		"no-shift lower": {sched: PartList{Text: "50,55", List: []int{50, 55}}, dt: 30, want: 50, wantShift: 0},
		"shift higher":   {sched: PartList{Text: "10,13,17", List: []int{10, 13, 17}}, dt: 23, want: 10, wantShift: 1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, shift := test.sched.compareTime(test.dt)
			if !reflect.DeepEqual(test.want, got) && !reflect.DeepEqual(test.wantShift, shift) {
				t.Fatalf("Expected: %#v, got: %#v", test.want, got)
			}
		})
	}
}
