package schedule

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCheckSchedule(t *testing.T) {

	partAll := PartAny{Text: "*"}
	// intervals
	int05 := PartInterval{Text: "0-5", Min: 0, Max: 5}
	int65 := PartInterval{Text: "6-5", Min: 6, Max: 5}
	int25 := PartInterval{Text: "0-25", Min: 0, Max: 25}
	int513 := PartInterval{Text: "5-13", Min: 5, Max: 13}
	// lists
	list50 := PartList{Text: "0,50", List: []int{0, 50}}
	list42 := PartList{Text: "4,2", List: []int{4, 2}}
	listSingle := PartList{Text: "23", List: []int{23}}
	list09 := PartList{Text: "0,9", List: []int{0, 9}}
	list513 := PartList{Text: "5,6,10,15", List: []int{5, 6, 10, 15}}

	tests := map[string]struct {
		part    Part
		partLim [2]int
		want    error
	}{
		"any: no error (ever)":           {part: partAll, partLim: [2]int{0, 59}, want: nil},
		"interval: correct (minute)":     {part: int05, partLim: [2]int{0, 59}, want: nil},
		"interval: min > max (weekDay)":  {part: int65, partLim: [2]int{0, 6}, want: fmt.Errorf("The ranges must be defined as 'min-max' with `min` <= `max`. Expects %v <= %v from string %s", int65.Min, int65.Max, int65.Text)},
		"interval: min lower (monthDay)": {part: int25, partLim: [2]int{1, 31}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 31, int25.Min, int25.Max, int25.Text)},
		"interval: max higher (month)":   {part: int513, partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 12, int513.Min, int513.Max, int513.Text)},
		"list: correct (minute)":         {part: list50, partLim: [2]int{0, 59}, want: nil},
		"list: min > max (weekDay)":      {part: list42, partLim: [2]int{0, 6}, want: nil},
		"list: single value (hour)":      {part: listSingle, partLim: [2]int{0, 23}, want: nil},
		"list: min lower (month)":        {part: list09, partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 12, 0, 9, list09.Text)},
		"list: max higher (month)":       {part: list513, partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 12, 5, 15, list513.Text)},
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
		Month:    PartList{Text: "1,2", List: []int{1, 2}},
	}

	listHours := Schedule{
		Second:   PartList{Text: "0", List: []int{0}},
		Minute:   PartList{Text: "0", List: []int{0}},
		Hour:     PartList{Text: "3,5,6", List: []int{3, 5, 6}},
		WeekDay:  every,
		MonthDay: PartList{Text: "31", List: []int{31}},
		Month:    every,
	}

	intervals := Schedule{
		Second:   PartInterval{Text: "55-58", Min: 55, Max: 58},
		Minute:   PartInterval{Text: "23-29", Min: 23, Max: 29},
		Hour:     PartInterval{Text: "3-6", Min: 3, Max: 6},
		WeekDay:  PartInterval{Text: "2-5", Min: 2, Max: 5},
		MonthDay: PartInterval{Text: "24-29", Min: 24, Max: 29},
		Month:    PartInterval{Text: "1-3", Min: 1, Max: 3},
	}

	tests := map[string]struct {
		sch  string
		want Schedule
		err  error
	}{
		"every second":                   {sch: "* * * * * *", want: everySecond, err: nil},
		"range minutes on Monday":        {sch: "0 2-5 * 0 * *", want: minutesMonday, err: nil},
		"specific time on 5th January ":  {sch: "0 30 12 * 5 1,2", want: specific, err: nil},
		"list hours every 31th monthDay": {sch: "0 0 3,5,6 * 31 *", want: listHours, err: nil},
		"error monthDay too high":        {sch: "0 0 12 * 32 *", want: Schedule{}, err: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 31, 32, 32, "32")},
		"error too many fields":          {sch: "0 0 0 0 0 0 *", want: Schedule{}, err: fmt.Errorf("Incorrect number of fields, expected 6 got %v. Fields are separated by a space and the whitespace can't be used for any other purpose", 7)},
		"error wrong range":              {sch: "0 0 12-18-10 0 0 *", want: Schedule{}, err: fmt.Errorf("Incorrect format of range. Expected 2 values separated by `-`, got %v", 3)},
		"error single non-convertible":   {sch: "a b c d e *", want: Schedule{}, err: fmt.Errorf("Unable to parse part of schedule: %s", "a")},
		"error non-convertible range":    {sch: "0 0 0 0 12-18a *", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "18a")},
		"error non-convertible list":     {sch: "0 0 0 0 12,1a *", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "1a")},
		"error list and range":           {sch: "0 11-15,16 0 0 0 *", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "15,16")},
		"error range and list":           {sch: "0 9,17-20 0 0 0 *", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "9,17")},
		"only intervals":                 {sch: "55-58 23-29 3-6 2-5 24-29 1-3", want: intervals, err: nil},
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

func TestCompareTime(t *testing.T) {
	tests := map[string]struct {
		sched     Part
		dt        int
		want      int
		wantShift int
	}{
		"any: no-shift zero":        {sched: PartAny{Text: "*"}, dt: 0, want: 0, wantShift: 0},
		"any: no-shift fifty-nine":  {sched: PartAny{Text: "*"}, dt: 59, want: 59, wantShift: 0},
		"interval: no-shift in set": {sched: PartInterval{Text: "0-3", Min: 0, Max: 3}, dt: 2, want: 2, wantShift: 0},
		"interval: no-shift lower":  {sched: PartInterval{Text: "5-10", Min: 5, Max: 10}, dt: 2, want: 5, wantShift: 0},
		"interval: shift higher":    {sched: PartInterval{Text: "4-9", Min: 4, Max: 9}, dt: 10, want: 4, wantShift: 1},
		"list: no-shift exact":      {sched: PartList{Text: "27", List: []int{27}}, dt: 27, want: 27, wantShift: 0},
		"list: no-shift lower":      {sched: PartList{Text: "50,55", List: []int{50, 55}}, dt: 30, want: 50, wantShift: 0},
		"list: shift higher":        {sched: PartList{Text: "10,13,17", List: []int{10, 13, 17}}, dt: 23, want: 10, wantShift: 1},
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

func TestNext(t *testing.T) {

	every := PartAny{Text: "*"}

	everySecond := Schedule{
		Second:   every,
		Minute:   every,
		Hour:     every,
		WeekDay:  every,
		MonthDay: every,
		Month:    every,
	}

	specificMDHMS := Schedule{
		Second:   PartList{Text: "0", List: []int{0}},
		Minute:   PartList{Text: "30", List: []int{30}},
		Hour:     PartList{Text: "12", List: []int{12}},
		WeekDay:  every,
		MonthDay: PartList{Text: "5", List: []int{5}},
		Month:    PartList{Text: "1", List: []int{1}},
	}

	listHMS := Schedule{
		Second:   PartList{Text: "55", List: []int{55}},
		Minute:   PartInterval{Text: "46-50", Min: 46, Max: 50},
		Hour:     PartInterval{Text: "16-20", Min: 16, Max: 20},
		WeekDay:  every,
		MonthDay: every,
		Month:    every,
	}

	specificWDMD := Schedule{
		Second:   PartList{Text: "1", List: []int{1}},
		Minute:   PartList{Text: "1", List: []int{1}},
		Hour:     PartList{Text: "1", List: []int{1}},
		WeekDay:  PartList{Text: "4", List: []int{4}},
		MonthDay: PartList{Text: "10", List: []int{10}},
		Month:    every,
	}

	listWDMD := Schedule{
		Second:   PartList{Text: "50", List: []int{50}},
		Minute:   PartList{Text: "26", List: []int{26}},
		Hour:     PartList{Text: "14", List: []int{14}},
		WeekDay:  PartList{Text: "2", List: []int{2}},
		MonthDay: PartList{Text: "10,11", List: []int{10, 11}},
		Month:    every,
	}

	intWDintMD := Schedule{
		Second:   PartList{Text: "2", List: []int{2}},
		Minute:   PartList{Text: "3", List: []int{3}},
		Hour:     PartList{Text: "4", List: []int{4}},
		WeekDay:  PartInterval{Text: "2-4", Min: 2, Max: 4},
		MonthDay: PartInterval{Text: "20-21", Min: 20, Max: 21},
		Month:    PartInterval{Text: "2-4", Min: 2, Max: 4},
	}

	tuesday2 := Schedule{
		Second:   PartList{Text: "5", List: []int{5}},
		Minute:   PartList{Text: "33", List: []int{33}},
		Hour:     PartList{Text: "15", List: []int{15}},
		WeekDay:  PartList{Text: "2", List: []int{2}},
		MonthDay: PartInterval{Text: "1-2", Min: 1, Max: 2},
		Month:    PartInterval{Text: "3-8", Min: 3, Max: 8},
	}

	tests := map[string]struct {
		sched Schedule
		after time.Time
		want  time.Time
	}{
		"this second":        {sched: everySecond, after: time.Date(2019, time.Month(10), 7, 23, 20, 0, 0, time.Local), want: time.Date(2019, time.Month(10), 7, 23, 20, 0, 0, time.Local)},
		"next year":          {sched: specificMDHMS, after: time.Date(2019, time.Month(10), 7, 23, 20, 0, 0, time.Local), want: time.Date(2020, time.Month(1), 5, 12, 30, 0, 0, time.Local)},
		"in 10 minutes":      {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 16, 35, 0, 0, time.Local), want: time.Date(2019, time.Month(3), 25, 16, 46, 55, 0, time.Local)},
		"in an hour":         {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 16, 51, 0, 0, time.Local), want: time.Date(2019, time.Month(3), 25, 17, 46, 55, 0, time.Local)},
		"in 5 seconds":       {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 20, 48, 53, 0, time.Local), want: time.Date(2019, time.Month(3), 25, 20, 48, 55, 0, time.Local)},
		"tomorrow":           {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 21, 8, 6, 0, time.Local), want: time.Date(2019, time.Month(3), 26, 16, 46, 55, 0, time.Local)},
		"no shift wDay":      {sched: specificWDMD, after: time.Date(2019, time.Month(10), 7, 1, 1, 1, 0, time.Local), want: time.Date(2019, time.Month(10), 10, 1, 1, 1, 0, time.Local)},
		"shift wDay list":    {sched: listWDMD, after: time.Date(2019, time.Month(10), 7, 19, 56, 38, 0, time.Local), want: time.Date(2019, time.Month(12), 10, 14, 26, 50, 0, time.Local)},
		"shift intervals":    {sched: intWDintMD, after: time.Date(2019, time.Month(10), 27, 22, 39, 16, 55, time.Local), want: time.Date(2020, time.Month(2), 20, 4, 3, 2, 0, time.Local)},
		"tuesday the second": {sched: tuesday2, after: time.Date(2019, time.Month(10), 1, 0, 0, 0, 1, time.Local), want: time.Date(2020, time.Month(6), 2, 15, 33, 5, 0, time.Local)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.sched.Next(test.after)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %v, got: %v", test.want, got)
			}
		})
	}

}
