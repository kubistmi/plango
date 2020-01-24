package schedule

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kubistmi/plango/utils"
)

func wrapMakeRange(min, max int) []int {
	res, _ := utils.MakeRange(min, max)
	return res
}

var every = partAny{Text: "*"}

func TestCheckSchedule(t *testing.T) {

	// interval specification
	int05 := partList{Text: "0-5", List: wrapMakeRange(0, 5)}
	int25 := partList{Text: "0-25", List: wrapMakeRange(0, 25)}
	int513 := partList{Text: "5-13", List: wrapMakeRange(5, 13)}
	// list specification
	list50 := partList{Text: "0,50", List: []int{0, 50}}
	list42 := partList{Text: "4,2", List: []int{4, 2}}
	listSingle := partList{Text: "23", List: []int{23}}
	list09 := partList{Text: "0,9", List: []int{0, 9}}
	list513 := partList{Text: "5,6,10,15", List: []int{5, 6, 10, 15}}

	tests := map[string]struct {
		part    part
		partLim [2]int
		want    error
	}{
		"interval: correct (minute)":     {part: int05, partLim: [2]int{0, 59}, want: nil},
		"any: no error (ever)":           {part: every, partLim: [2]int{0, 59}, want: nil},
		"interval: min lower (monthDay)": {part: int25, partLim: [2]int{1, 31}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 31, int25.min(), int25.max(), int25.Text)},
		"interval: max higher (month)":   {part: int513, partLim: [2]int{1, 12}, want: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 12, int513.min(), int513.max(), int513.Text)},
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

	everySecond := Schedule{
		Second:   every,
		Minute:   every,
		Hour:     every,
		MonthDay: every,
		Month:    every,
		WeekDay:  every,
	}

	minutesMonday := Schedule{
		Second:   partList{Text: "0", List: []int{0}},
		Minute:   partList{Text: "2-5", List: wrapMakeRange(2, 5)},
		Hour:     every,
		MonthDay: every,
		Month:    every,
		WeekDay:  partList{Text: "0", List: []int{0}},
	}

	specific := Schedule{
		Second:   partList{Text: "0", List: []int{0}},
		Minute:   partList{Text: "30", List: []int{30}},
		Hour:     partList{Text: "12", List: []int{12}},
		MonthDay: partList{Text: "5", List: []int{5}},
		Month:    partList{Text: "1,2", List: []int{1, 2}},
		WeekDay:  every,
	}

	listHours := Schedule{
		Second:   partList{Text: "0", List: []int{0}},
		Minute:   partList{Text: "0", List: []int{0}},
		Hour:     partList{Text: "3,5,6", List: []int{3, 5, 6}},
		MonthDay: partList{Text: "31", List: []int{31}},
		Month:    every,
		WeekDay:  every,
	}

	intervals := Schedule{
		Second:   partList{Text: "55-58", List: wrapMakeRange(55, 58)},
		Minute:   partList{Text: "23-29", List: wrapMakeRange(23, 29)},
		Hour:     partList{Text: "3-6", List: wrapMakeRange(3, 6)},
		MonthDay: partList{Text: "24-29", List: wrapMakeRange(24, 29)},
		Month:    partList{Text: "1-3", List: wrapMakeRange(1, 3)},
		WeekDay:  partList{Text: "5-2", List: wrapMakeRange(2, 5)},
	}

	tests := map[string]struct {
		sch  string
		want Schedule
		err  error
	}{
		"every second":                   {sch: "* * * * * *", want: everySecond, err: nil},
		"range minutes on Monday":        {sch: "0 2-5 * * * 0", want: minutesMonday, err: nil},
		"specific time on 5th January ":  {sch: "0 30 12 5 1,2 *", want: specific, err: nil},
		"list hours every 31th monthDay": {sch: "0 0 3,5,6 31 * *", want: listHours, err: nil},
		"error monthDay too high":        {sch: "0 0 12 32 * *", want: Schedule{}, err: fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s", 1, 31, 32, 32, "32")},
		"error too many fields":          {sch: "0 0 0 0 0 * 0", want: Schedule{}, err: fmt.Errorf("Incorrect number of fields, expected 6 got %v. Fields are separated by a space and the whitespace can't be used for any other purpose", 7)},
		"error wrong range":              {sch: "0 0 12-18-10 0 * 0", want: Schedule{}, err: fmt.Errorf("Incorrect format of range. Expected 2 values separated by `-`, got %v", 3)},
		"error single non-convertible":   {sch: "a b c e * d", want: Schedule{}, err: fmt.Errorf("Unable to parse part of schedule: %s", "a")},
		"error non-convertible range":    {sch: "0 0 0 12-18a * 0", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "18a")},
		"error non-convertible list":     {sch: "0 0 0 12,1a * 0", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "1a")},
		"error list and range":           {sch: "0 11-15,16 0 0 * 0", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "15,16")},
		"error range and list":           {sch: "0 9,17-20 0 0 * 0", want: Schedule{}, err: fmt.Errorf("Unable to convert %s to an integer", "9,17")},
		"only intervals":                 {sch: "55-58 23-29 3-6 24-29 1-3 5-2", want: intervals, err: nil},
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
		sched     part
		dt        int
		want      int
		wantShift int
	}{
		"any: no-shift zero":        {sched: partAny{Text: "*"}, dt: 0, want: 0, wantShift: 0},
		"any: no-shift fifty-nine":  {sched: partAny{Text: "*"}, dt: 59, want: 59, wantShift: 0},
		"interval: no-shift in set": {sched: partList{Text: "0-3", List: wrapMakeRange(0, 3)}, dt: 2, want: 2, wantShift: 0},
		"interval: no-shift lower":  {sched: partList{Text: "5-10", List: wrapMakeRange(5, 10)}, dt: 2, want: 5, wantShift: 0},
		"interval: shift higher":    {sched: partList{Text: "4-9", List: wrapMakeRange(4, 9)}, dt: 10, want: 4, wantShift: 1},
		"list: no-shift exact":      {sched: partList{Text: "27", List: []int{27}}, dt: 27, want: 27, wantShift: 0},
		"list: no-shift lower":      {sched: partList{Text: "50,55", List: []int{50, 55}}, dt: 30, want: 50, wantShift: 0},
		"list: shift higher":        {sched: partList{Text: "10,13,17", List: []int{10, 13, 17}}, dt: 23, want: 10, wantShift: 1},
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

	every := partAny{Text: "*"}

	everySecond := Schedule{
		Second:   every,
		Minute:   every,
		Hour:     every,
		WeekDay:  every,
		MonthDay: every,
		Month:    every,
	}

	specificMDHMS := Schedule{
		Second:   partList{Text: "0", List: []int{0}},
		Minute:   partList{Text: "30", List: []int{30}},
		Hour:     partList{Text: "12", List: []int{12}},
		WeekDay:  every,
		MonthDay: partList{Text: "5", List: []int{5}},
		Month:    partList{Text: "1", List: []int{1}},
	}

	listHMS := Schedule{
		Second:   partList{Text: "55", List: []int{55}},
		Minute:   partList{Text: "46-50", List: wrapMakeRange(46, 50)},
		Hour:     partList{Text: "16-20", List: wrapMakeRange(16, 20)},
		WeekDay:  every,
		MonthDay: every,
		Month:    every,
	}

	specificWDMD := Schedule{
		Second:   partList{Text: "1", List: []int{1}},
		Minute:   partList{Text: "1", List: []int{1}},
		Hour:     partList{Text: "1", List: []int{1}},
		WeekDay:  partList{Text: "4", List: []int{4}},
		MonthDay: partList{Text: "10", List: []int{10}},
		Month:    every,
	}

	listWDMD := Schedule{
		Second:   partList{Text: "50", List: []int{50}},
		Minute:   partList{Text: "26", List: []int{26}},
		Hour:     partList{Text: "14", List: []int{14}},
		WeekDay:  partList{Text: "2", List: []int{2}},
		MonthDay: partList{Text: "10,11", List: []int{10, 11}},
		Month:    every,
	}

	intWDintMD := Schedule{
		Second:   partList{Text: "2", List: []int{2}},
		Minute:   partList{Text: "3", List: []int{3}},
		Hour:     partList{Text: "4", List: []int{4}},
		WeekDay:  partList{Text: "2-4", List: wrapMakeRange(2, 4)},
		MonthDay: partList{Text: "20-21", List: wrapMakeRange(20, 21)},
		Month:    partList{Text: "2-4", List: wrapMakeRange(2, 4)},
	}

	tuesday2 := Schedule{
		Second:   partList{Text: "5", List: []int{5}},
		Minute:   partList{Text: "33", List: []int{33}},
		Hour:     partList{Text: "15", List: []int{15}},
		WeekDay:  partList{Text: "2", List: []int{2}},
		MonthDay: partList{Text: "1-2", List: wrapMakeRange(1, 2)},
		Month:    partList{Text: "3-8", List: wrapMakeRange(3, 8)},
	}

	nextmonth31 := Schedule{
		Second:   partList{Text: "5", List: []int{5}},
		Minute:   partList{Text: "33", List: []int{33}},
		Hour:     partList{Text: "15", List: []int{15}},
		WeekDay:  every,
		MonthDay: partList{Text: "31", List: []int{31}},
		Month:    partList{Text: "11,12", List: []int{11, 12}},
	}

	tests := map[string]struct {
		sched Schedule
		after time.Time
		want  time.Time
		err   error
	}{
		"this second":                    {sched: everySecond, after: time.Date(2019, time.Month(10), 7, 23, 20, 0, 0, time.Local), want: time.Date(2019, time.Month(10), 7, 23, 20, 0, 0, time.Local)},
		"next year":                      {sched: specificMDHMS, after: time.Date(2019, time.Month(10), 7, 23, 20, 0, 0, time.Local), want: time.Date(2020, time.Month(1), 5, 12, 30, 0, 0, time.Local)},
		"in 10 minutes":                  {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 16, 35, 0, 0, time.Local), want: time.Date(2019, time.Month(3), 25, 16, 46, 55, 0, time.Local)},
		"in an hour":                     {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 16, 51, 0, 0, time.Local), want: time.Date(2019, time.Month(3), 25, 17, 46, 55, 0, time.Local)},
		"in 5 seconds":                   {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 20, 48, 53, 0, time.Local), want: time.Date(2019, time.Month(3), 25, 20, 48, 55, 0, time.Local)},
		"tomorrow":                       {sched: listHMS, after: time.Date(2019, time.Month(3), 25, 21, 8, 6, 0, time.Local), want: time.Date(2019, time.Month(3), 26, 16, 46, 55, 0, time.Local)},
		"no shift wDay":                  {sched: specificWDMD, after: time.Date(2019, time.Month(10), 7, 1, 1, 1, 0, time.Local), want: time.Date(2019, time.Month(10), 10, 1, 1, 1, 0, time.Local)},
		"shift wDay list":                {sched: listWDMD, after: time.Date(2019, time.Month(10), 7, 19, 56, 38, 0, time.Local), want: time.Date(2019, time.Month(12), 10, 14, 26, 50, 0, time.Local)},
		"shift intervals":                {sched: intWDintMD, after: time.Date(2019, time.Month(10), 27, 22, 39, 16, 55, time.Local), want: time.Date(2020, time.Month(2), 20, 4, 3, 2, 0, time.Local)},
		"tuesday the second":             {sched: tuesday2, after: time.Date(2019, time.Month(10), 1, 0, 0, 0, 1, time.Local), want: time.Date(2020, time.Month(6), 2, 15, 33, 5, 0, time.Local)},
		"skipping the month with day 31": {sched: nextmonth31, after: time.Date(2019, time.Month(10), 30, 0, 0, 0, 1, time.Local), want: time.Date(2019, time.Month(12), 31, 0, 0, 1, 0, time.Local)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := test.sched.Next(test.after)
			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("Expected: %v, got: %v", test.want, got)
			}
			if !reflect.DeepEqual(test.err, err) {
				t.Fatalf("Expected: %#v, got: %#v", test.err, err)
			}
		})
	}

}
