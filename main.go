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

// finMin ...
func findMin(vec []int) int {
	var min int
	for ix, val := range vec {
		if ix == 0 || val < min {
			min = val
		}
	}
	return min
}

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
var PartLimits = map[string][2]int{
	"Second":   [2]int{0, 59},
	"Minute":   [2]int{0, 59},
	"Hour":     [2]int{0, 23},
	"WeekDay":  [2]int{0, 6},
	"MonthDay": [2]int{1, 31},
	"Month":    [2]int{1, 12},
}

// Schedule ...
type Schedule struct {
	Second, Minute, Hour, WeekDay, MonthDay, Month ParsedPart
}

// ParsedPart ...
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
*/

// CompareTime ...
func CompareTime(sched ParsedPart, dt int) (int, int) {
	if sched.Any {
		return dt, 0
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
		if val >= dt {
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

//NewSchedule("{s:,m:,h:,wd:,md:,m:}")

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

// ParseSchedule finds the time that the Job should be executed
func ParseSchedule(Schedule string) time.Time {
	return time.Now()
}

func main() {
	fmt.Println(PartLimits)
}
