package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// NumSchedules defines the number of schedules that stored
	NumSchedules = 5
)

// findMin ...
// TODO: packafe utils
func findMin(vec []int) int {
	var min int
	for ix, val := range vec {
		if ix == 0 || val < min {
			min = val
		}
	}
	return min
}

// findMax ...
// TODO: packafe utils
func findMax(vec []int) int {
	var max int
	for ix, val := range vec {
		if ix == 0 || val > max {
			max = val
		}
	}
	return max
}

// makeRange ...
// TODO: package utils
func makeRange(min, max int) ([]int, error) {
	if min >= max {
		return []int{}, fmt.Errorf("The parameter `max` must be strictly larger than parameter `min`")
	}

	vec := make([]int, max-min+1)
	for ix := range vec {
		vec[ix] = min + ix
	}
	return vec, nil
}

//findUnique
// TODO: package utils
func findUnique(vec []int) []int {
	unique := make([]int, 0, len(vec))
	mapper := make(map[int]bool)

	for _, val := range vec {
		if _, ok := mapper[val]; !ok {
			mapper[val] = true
			unique = append(unique, val)
		}
	}
	return unique
}

// PartLimits ...
// TODO: file schedule
var PartLimits = map[string][2]int{
	"second":   [2]int{0, 59},
	"minute":   [2]int{0, 59},
	"hour":     [2]int{0, 23},
	"weekDay":  [2]int{0, 6},
	"monthDay": [2]int{1, 31},
	"month":    [2]int{1, 12},
}

// PartOrder ...
// TODO: file schedule
var PartOrder = []string{"second", "minute", "hour", "weekDay", "monthDay", "month"}

// SchedulePart ...
// TODO: file schedule
type SchedulePart interface {
	checkPart(partLim [2]int) error
	compareTime(timepart int) (int, int)
	getOrigin() string
}

// Schedule ...
// TODO: file schedule
type Schedule struct{ Second, Minute, Hour, WeekDay, MonthDay, Month SchedulePart }

// PartAny defines schedule based on the string "*". This definition will trigger on every occurence it can.
// E.g. using * in monthDays field means the job will be run every day of the month (with regard to other definitions, such as weekDay).
type PartAny struct {
	Text string
}

// PartInterval defines schedule based on the string "x-y". This definition will trigger on every occurence in this interval.
// E.g. using 4-6 in hours field means the job will be run at hours 4, 5 and 6 (with regard to other definitions).
type PartInterval struct {
	Min, Max int
	Text     string
}

// PartList defines schedule based on the string "x,y,y". This definition will trigger on every occurence listed in the definition.
// E.g. using 4,6,20 in minutes field means the job will be run at minutes 4, 6 and 20 (with regard to other definitions).
type PartList struct {
	List []int
	Text string
}

// ---------------------- GET ORIGIN -------------------------------------------------------------------
func (sp PartAny) getOrigin() string {
	return sp.Text
}

func (sp PartInterval) getOrigin() string {
	return sp.Text
}

func (sp PartList) getOrigin() string {
	return sp.Text
}

// ---------------------- COMPARE TIME -----------------------------------------------------------------

// compareTime ...
func (sp PartAny) compareTime(timepart int) (int, int) {
	// timepart originates from time.Time so it's valid time value => safe to return
	return timepart, 0
}

func (sp PartInterval) compareTime(timepart int) (int, int) {
	var vec []int
	var err error

	vec, err = makeRange(sp.Min, sp.Max)
	if err != nil {
		log.Fatal(err)
	}

	// find time value higher than or equal current time
	// e.g. schedule = "0 0 5-7 * * *" ; current = 2019-10-12 06:00:00
	//      return 2019-10-12 06:00:00
	for _, val := range vec {
		if val >= timepart {
			return val, 0
		}
	}

	// if no such value, choose the next (the smallest one) and shift the time
	// e.g. schedule = "0 0 2-5 * * *" ; current = 2019-10-12 06:00:00
	//      return 2019-10-13 02:00:00
	return sp.Min, 1
}

func (sp PartList) compareTime(timepart int) (int, int) {
	// find time value higher than or equal current time
	// e.g. schedule = "0 0 5,6 * * *" ; current = 2019-10-12 06:00:00
	//      return 2019-10-12 06:00:00
	for _, val := range sp.List {
		if val >= timepart {
			return val, 0
		}
	}

	// if no such value, choose the next (the smallest one) and shift the time
	// e.g. schedule = "0 0 2,3,4 * * *" ; current = 2019-10-12 06:00:00
	//      return 2019-10-13 02:00:00
	return findMin(sp.List), 1
}

// ---------------------- CHECK PART ---------------------------------------------------------------

// checkPart ...
// TODO: file schedule
func (sp PartAny) checkPart(partLim [2]int) error {
	return nil
}

func (sp PartInterval) checkPart(partLim [2]int) error {
	if sp.Min > sp.Max {
		return fmt.Errorf("The ranges must be defined as 'min-max' with `min` <= `max`. Expects %v <= %v from string %s",
			sp.Min, sp.Max, sp.Text)
	}
	if !(sp.Min >= partLim[0] && sp.Min <= partLim[1] && sp.Max >= partLim[0] && sp.Max <= partLim[1]) {
		return fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s",
			partLim[0], partLim[1], sp.Min, sp.Max, sp.Text)
	}
	return nil
}

func (sp PartList) checkPart(partLim [2]int) error {
	min := findMin(sp.List)
	max := findMax(sp.List)

	if min > max {
		return fmt.Errorf("The ranges must be defined as 'min-max' with `min` <= `max`. Expects %v <= %v from string %s",
			min, max, sp.Text)
	}
	if !(min >= partLim[0] && min <= partLim[1] && max >= partLim[0] && max <= partLim[1]) {
		return fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v from string %s",
			partLim[0], partLim[1], min, max, sp.Text)
	}
	return nil

}

// ParseSchedule ...
func ParseSchedule(schedule string) (Schedule, error) {

	res := make(map[string]SchedulePart)

	parts := strings.Split(schedule, " ")

	if lenParts := len(parts); lenParts != len(PartOrder) {
		return Schedule{},
			fmt.Errorf("Incorrect number of fields, expected 6 got %v. Fields are separated by a space and the whitespace can't be used for any other purpose", lenParts)
	}

	// process part-by-part
	for ix, p := range parts {
		var part SchedulePart
		var err error

		partType := PartOrder[ix]
		partLim := PartLimits[partType]

		switch {
		case strings.Contains(p, "*"):
			part = PartAny{Text: p}

		case strings.Contains(p, "-"):
			lims := strings.Split(p, "-")

			if lenLim := len(lims); lenLim != 2 {
				return Schedule{},
					fmt.Errorf("Incorrect format of range. Expected 2 values separated by `-`, got %v", lenLim)
			}

			limsI := make([]int, len(lims))
			for ix, val := range lims {
				limsI[ix], err = strconv.Atoi(val)
				if err != nil {
					return Schedule{},
						fmt.Errorf("Unable to convert %s to an integer", val)
				}
			}

			min := findMin(limsI)
			max := findMax(limsI)

			part = PartInterval{
				Text: p,
				Min:  min,
				Max:  max,
			}

		case strings.Contains(p, ","):
			list := strings.Split(p, ",")

			listI := make([]int, len(list))

			for ix, val := range list {
				listI[ix], err = strconv.Atoi(val)
				if err != nil {
					return Schedule{}, fmt.Errorf("Unable to convert %s to an integer", val)
				}
			}

			// sort and keep unique only
			sort.Ints(listI)
			part = PartList{Text: p, List: findUnique(listI)}

		default:
			pI, err := strconv.Atoi(p)
			if err != nil {
				return Schedule{}, fmt.Errorf("Unable to parse part of schedule: %s", p)
			}

			part = PartList{Text: p, List: []int{pI}}
		}
		err = part.checkPart(partLim)
		if err != nil {
			return Schedule{}, err
		}

		res[partType] = part
	}

	return Schedule{
		Second:   res["second"],
		Minute:   res["minute"],
		Hour:     res["hour"],
		WeekDay:  res["weekDay"],
		MonthDay: res["monthDay"],
		Month:    res["month"],
	}, nil

}

// Next ...
func (s Schedule) Next(After time.Time) time.Time {
	var nxtSecond, nxtMinute, nxtHour, nxtMday, nxtMonth, nxtYear int //nxtWday,
	var shift int

	next := After

	nxtSecond, shift = s.Second.compareTime(next.Second())
	next = next.Add(time.Minute * time.Duration(shift))

	nxtMinute, shift = s.Minute.compareTime(next.Minute())
	next = next.Add(time.Hour * time.Duration(shift))

	nxtHour, shift = s.Hour.compareTime(next.Hour())
	next = next.AddDate(0, 0, shift)

	// day is a little bit more fun
	switch s.WeekDay.(type) {

	case PartAny:
		nxtMday, shift = s.MonthDay.compareTime(next.Day())
		next = next.AddDate(0, shift, 0)

		nxtMonth, shift = s.Month.compareTime(int(next.Month()))
		nxtYear = next.Year() + shift
		return time.Date(nxtYear, time.Month(nxtMonth), nxtMday, nxtHour, nxtMinute, nxtSecond, 0, time.Local)

	default:
		var wdList []int

		switch v := s.WeekDay.(type) {
		case PartInterval:
			wdList, _ = makeRange(v.Min, v.Max)

		case PartList:
			wdList = v.List
		}

		iter := len(wdList) * 53

		next = time.Date(next.Year(), next.Month(), next.Day(), nxtHour, nxtMinute, nxtSecond, 0, time.Local)
		var wdMday, wdNext, wdShift int

		for i := 0; i < iter; i++ {
			wdNext, wdShift = s.WeekDay.compareTime(int(next.Weekday()))

			if wdShift == 1 {
				next = next.AddDate(0, 0, 7-int(next.Weekday()))
			}
			next = next.AddDate(0, 0, wdNext-int(next.Weekday()))

			wdMday, wdShift = s.MonthDay.compareTime(next.Day())

			if wdShift == 0 && wdMday == int(next.Day()) {
				return time.Date(next.Year(), next.Month(), next.Day(), nxtHour, nxtMinute, nxtSecond, 0, time.Local)
			}
			if wdShift == 1 {
				next = next.AddDate(0, 1, wdMday-int(next.Day()))

			} else {
				next = next.AddDate(0, 0, wdMday-int(next.Day()))
			}
		}
	}
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
}

// Job contains the task definition
type Job struct {
	ID   string
	Name string
	Schedule
	Active  bool
	Command string
	Args    []string
	Config  map[string]string
}

// Run defines the singular execution of the Job
type Run struct {
	ID        string
	JobID     string
	StartTime time.Time
	EndTime   time.Time
	Status    string
	Trigger   string
}

// CreateJob prepares a new Job definition
func CreateJob(Name string, Plan Schedule, Command string, Args []string, Config map[string]string) Job {
	// TODO: implement JobID collection
	NewJob := Job{
		ID:       "1",
		Name:     Name,
		Schedule: Plan,
		Command:  Command,
		Args:     Args,
		Config:   Config,
	}

	// TODO: implement SQL upload

	// TODO: implement schedule preparation

	return (NewJob)

}

func main() {
	fmt.Println("Build succesfull!")
}
