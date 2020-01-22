// Package schedule ...
package schedule

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kubistmi/plango/utils"
)

// PartLimits ...
var partLimits = map[string][2]int{
	"second":   [2]int{0, 59},
	"minute":   [2]int{0, 59},
	"hour":     [2]int{0, 23},
	"monthDay": [2]int{1, 31},
	"month":    [2]int{1, 12},
	"weekDay":  [2]int{0, 6},
}

// PartOrder ...
var partOrder = []string{"second", "minute", "hour", "monthDay", "month", "weekDay"}

// Part ...
type part interface {
	checkPart(partLim [2]int) error
	compareTime(timepart int) (int, int)
	getOrigin() string
}

// Schedule ...
type Schedule struct{ Second, Minute, Hour, MonthDay, Month, WeekDay part }

// PartAny defines schedule based on the string "*". This definition will trigger on every occurence it can.
// E.g. using * in monthDays field means the job will be run every day of the month (with regard to other definitions, such as weekDay).
type partAny struct {
	Text string
}

// PartInterval defines schedule based on the string "x-y". This definition will trigger on every occurence in this interval.
// E.g. using 4-6 in hours field means the job will be run at hours 4, 5 and 6 (with regard to other definitions).
type partInterval struct {
	Min, Max int
	Text     string
}

// PartList defines schedule based on the string "x,y,y". This definition will trigger on every occurence listed in the definition.
// E.g. using 4,6,20 in minutes field means the job will be run at minutes 4, 6 and 20 (with regard to other definitions).
type partList struct {
	List []int
	Text string
}

// ---------------------- GET ORIGIN -------------------------------------------------------------------
func (sp partAny) getOrigin() string {
	return sp.Text
}

func (sp partInterval) getOrigin() string {
	return sp.Text
}

func (sp partList) getOrigin() string {
	return sp.Text
}

// ---------------------- COMPARE TIME -----------------------------------------------------------------

// compareTime ...
func (sp partAny) compareTime(timepart int) (int, int) {
	// timepart originates from time.Time so it's valid time value => safe to return
	return timepart, 0
}

// compareTime ...
func (sp partInterval) compareTime(timepart int) (int, int) {
	var vec []int
	var err error

	// TODO: how to properly implement this?
	vec, err = utils.MakeRange(sp.Min, sp.Max)
	if err != nil {
		log.Fatalf(
			`Ok, so while comparing the time from the input %v (type Interval), the method was unable to create the slice defined by the bounds. \n
			This was supposed to be handled while parsing the schedule and therefore should NEVER happen. Congratulations, I am so sorry, and please submit an issue at https://github.com/kubistmi/plango. \n
			What you need to do now is to remove this schedule and resubmit it. Perhaps try different format. This is uncharted territory, so I really don't know.`,
			sp.Text)
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

// compareTime ...
func (sp partList) compareTime(timepart int) (int, int) {
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
	return utils.FindMin(sp.List), 1
}

// ---------------------- CHECK PART ---------------------------------------------------------------
// checkPart ...
func (sp partAny) checkPart(partLim [2]int) error {
	return nil
}

// checkPart ...
func (sp partInterval) checkPart(partLim [2]int) error {
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

// checkPart ...
func (sp partList) checkPart(partLim [2]int) error {
	min := utils.FindMin(sp.List)
	max := utils.FindMax(sp.List)

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

	res := make(map[string]part)

	parts := strings.Split(schedule, " ")

	if lenParts := len(parts); lenParts != len(partOrder) {
		return Schedule{},
			fmt.Errorf("Incorrect number of fields, expected 6 got %v. Fields are separated by a space and the whitespace can't be used for any other purpose", lenParts)
	}

	// process part-by-part
	for ix, p := range parts {
		var part part
		var err error

		partType := partOrder[ix]
		partLim := partLimits[partType]

		switch {
		case strings.Contains(p, "*"):
			part = partAny{Text: p}

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

			min := utils.FindMin(limsI)
			max := utils.FindMax(limsI)

			part = partInterval{
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
			part = partList{Text: p, List: utils.FindUnique(listI)}

		default:
			pI, err := strconv.Atoi(p)
			if err != nil {
				return Schedule{}, fmt.Errorf("Unable to parse part of schedule: %s", p)
			}

			part = partList{Text: p, List: []int{pI}}
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
		MonthDay: res["monthDay"],
		Month:    res["month"],
		WeekDay:  res["weekDay"],
	}, nil

}

// Next ...
func (s Schedule) Next(After time.Time) (time.Time, error) {
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

	case partAny:
		// the easy part, for non-specific weekDay, just go through the calendar as above
		nxtMday, shift = s.MonthDay.compareTime(next.Day())
		next = next.AddDate(0, shift, 0)

		nxtMonth, shift = s.Month.compareTime(int(next.Month()))
		nxtYear = next.Year() + shift
		return time.Date(nxtYear, time.Month(nxtMonth), nxtMday, nxtHour, nxtMinute, nxtSecond, 0, time.Local), nil

	default:
		// TODO: should this be a config variable?
		iter := 50

		// update the next timestamp using the found values of H:M:S
		next = time.Date(next.Year(), next.Month(), next.Day(), nxtHour, nxtMinute, nxtSecond, 0, time.Local)
		var wdMday, wdNext, wdShift int

		for i := 0; i < iter; i++ {
			// first, check whether the weekday is OK
			// shift by:
			//    - difference in days
			//    - difference in weekdays * 7
			wdNext, wdShift = s.WeekDay.compareTime(int(next.Weekday()))
			next = next.AddDate(0, 0, wdNext-int(next.Weekday())+wdShift*7)

			// then, check the monthDay and month
			wdMday, wdShift = s.MonthDay.compareTime(next.Day())
			wdMonth, wdMshift := s.Month.compareTime(int(next.Month()))

			// very greedy early exit, if the monthDay and month are OK, return
			if wdShift == 0 && wdMday == next.Day() && wdMonth == int(next.Month()) {
				return time.Date(next.Year(), next.Month(), next.Day(), nxtHour, nxtMinute, nxtSecond, 0, time.Local), nil
			}

			// if Schedule.Month too high, then jump to the next year and ...
			// else to the max(Schedule.Month, 1) and ...
			// ... lowest of the Schedule.MonthDay
			if wdMshift == 1 {
				next = next.AddDate(wdMshift, wdMonth-int(next.Month()), wdMday-int(next.Day()))
			} else {
				monthShift := []int{wdMonth - int(next.Month()), wdShift}
				next = next.AddDate(0, utils.FindMax(monthShift), wdMday-int(next.Day()))
			}
		}
	}
	// if the date cannot be found, not sure this is reachable
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), fmt.Errorf("unable to find the date satisfying the schedule")
}
