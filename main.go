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

// Schedule ...
// TODO: file schedule
type Schedule struct {
	Second, Minute, Hour, WeekDay, MonthDay, Month ParsedPart
}

// ParsedPart ...
// TODO: file schedule
type ParsedPart struct {
	Any  bool
	Min  int
	Max  int
	List []int
	Type string
}

// CheckSchedule ...
// TODO: file schedule
func CheckSchedule(min, max int, p string, partLim [2]int) error {

	if min >= max {
		return fmt.Errorf("The ranges must be defined as 'min-max' with `min` >= `max`. Expects %v >= %v from string %s",
			min, max, p)
	}

	if !(min >= partLim[0] && min <= partLim[1] && max >= partLim[0] && max <= partLim[1]) {
		return fmt.Errorf("The range is not compliant for this part of Schedule. Expects numbers between %v-%v, got %v-%v",
			partLim[0], partLim[1], min, max)
	}
	return nil
}

// ParseSchedule ...
func ParseSchedule(schedule string) (Schedule, error) {

	res := make(map[string]ParsedPart)

	parts := strings.Split(schedule, " ")

	if lenParts := len(parts); lenParts != len(PartOrder) {
		return Schedule{},
			fmt.Errorf("Incorrect number of fields, expected 6 got %v. Fields are separated by a space and the whitespace can't be used for any other purpose", lenParts)
	}

	// process part-by-part
	for ix, p := range parts {
		part := new(ParsedPart)
		var err error

		partType := PartOrder[ix]
		partLim := PartLimits[partType]

		switch {
		case strings.Contains(p, "*"):
			part.Any = true

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

			err = CheckSchedule(min, max, p, partLim)
			if err != nil {
				return Schedule{}, err
			}
			part.Min = min
			part.Max = max

		case strings.Contains(p, ","):
			list := strings.Split(p, ",")

			listI := make([]int, len(list))

			for ix, val := range list {
				listI[ix], err = strconv.Atoi(val)
				if err != nil {
					return Schedule{}, fmt.Errorf("Unable to convert %s to an integer", val)
				}
			}

			min := findMin(listI)
			max := findMax(listI)

			err = CheckSchedule(min, max, p, partLim)
			if err != nil {
				return Schedule{}, err
			}

			// sort and keep unique only
			sort.Ints(listI)
			part.List = findUnique(listI)
		default:
			pI, err := strconv.Atoi(p)
			if err != nil {
				return Schedule{}, fmt.Errorf("Unable to parse part of schedule: %s", p)
			}
			err = CheckSchedule(pI, partLim[1], p, partLim)
			if err != nil {
				return Schedule{}, err
			}
			part.List = []int{pI}
		}
		res[partType] = *part
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

// CompareTime ...
func CompareTime(sched ParsedPart, timepart int) (int, int) {
	if sched.Any {
		// timepart originates from time.Time so it's valid time value => safe to return
		return timepart, 0
	}

	var vec []int
	var err error
	if len(sched.List) > 0 {
		vec = sched.List
	} else if sched.Min+sched.Max != 0 {
		vec, err = makeRange(sched.Min, sched.Max)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, val := range vec {
		if val >= timepart {
			return val, 0
		}
	}
	return findMin(vec), 1

}

// Next ...
func (s Schedule) Next(After time.Time) time.Time {

	After.Second()
	After.Minute()
	After.Hour()
	After.Weekday()
	After.Day()
	After.Month()
	return time.Now()

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
	fmt.Println(PartLimits)
}
