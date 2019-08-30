package main

import (
	"fmt"
	"log"
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

/*
func CompareTime(sched ParsedPart, dt int,  int) int {
	if sched.Any {
		return dt
	} else if Len(sched.List) > 0 {
		for _, val := range sched.List {
			if val >= dt {
				return(val)
			}
		}
	} else if sched.Min + sched.Max != 0 {
	}
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
