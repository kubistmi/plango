package main

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/kubistmi/plango/schedule"
)

const (
	// NumSchedules defines the number of future runs that are planned and stored
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
	// TODO: is this needed?
	Config map[string]string
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
	//fmt.Println("Build succesfull!")
	badSchedule, err := schedule.ParseSchedule("* * * 31 10 1")
	fmt.Println(err)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	go http.ListenAndServe(":8080", nil)

	fmt.Println(badSchedule.Next(time.Now()))
	wait := time.Tick(50 * time.Second)
	<-wait
}
