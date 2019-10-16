package main

import (
	"fmt"
	"time"

	"github.com/kubistmi/plango/schedule"
)

const (
	// NumSchedules defines the number of schedules that stored
	NumSchedules = 5
)

// Job contains the task definition
type Job struct {
	ID       string
	Name     string
	Schedule schedule.Schedule
	Active   bool
	Command  string
	Args     []string
	Config   map[string]string
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
func CreateJob(Name string, Plan schedule.Schedule, Command string, Args []string, Config map[string]string) Job {
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
