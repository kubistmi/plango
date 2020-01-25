// Package schedule ...
package schedule

import (
	"fmt"
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
	min(date int) int
	isin(el int) bool
	checkPart(partLim [2]int) error
	compareTime(timepart int) (int, int)
	getOrigin() string
}

// Schedule ...
type Schedule struct{ Second, Minute, Hour, MonthDay, Month, WeekDay part }

// checkSchedule
func checkSchedule(sch Schedule) error {

	days := make([]int, 0, 31)
	months := make([]int, 0, 12)

	switch d := sch.MonthDay.(type) {
	case partList:
		days = d.List
	}

	switch m := sch.Month.(type) {
	case partList:
		months = m.List
	}

	if len(months)*len(days) == 0 {
		return nil
	}

	for _, m := range months {
		for _, d := range days {
			dt := time.Date(2000, time.Month(m), d, 0, 0, 0, 0, time.Local)
			if dt.Day() != d || int(dt.Month()) != m {
				return fmt.Errorf("Encoutered impossible schedule: month:%v - day:%v", m, d)
			}
		}
	}
	return nil
}

// PartAny defines schedule based on the string "*". This definition will trigger on every occurence it can.
// E.g. using * in monthDays field means the job will be run every day of the month (with regard to other definitions, such as weekDay).
type partAny struct {
	Text string
}

// PartList defines schedule based on the string "x,y,y". This definition will trigger on every occurence listed in the definition.
// E.g. using 4,6,20 in minutes field means the job will be run at minutes 4, 6 and 20 (with regard to other definitions).
type partList struct {
	List []int
	Text string
}

// ---------------------- ISIN -----------------------------------------------------------------------
func (sp partAny) isin(el int) bool {
	return true
}

func (sp partList) isin(el int) bool {
	return utils.IsIn(el, sp.List)
}

// ---------------------- MIN ---------------------------------------------------------------------------
func (sp partAny) min(date int) int {
	return 0 + date
}

func (sp partList) min(date int) int {
	return utils.FindMin(sp.List)
}

// ---------------------- GET ORIGIN -------------------------------------------------------------------
func (sp partAny) getOrigin() string {
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

			list, err := utils.MakeRange(min, max)
			//! not sure this is reachable
			if err != nil {
				return Schedule{},
					fmt.Errorf("Error when attempting to build a list. Expected min-max got %v-%v", min, max)
			}

			part = partList{
				Text: p,
				List: list,
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

	result := Schedule{
		Second:   res["second"],
		Minute:   res["minute"],
		Hour:     res["hour"],
		MonthDay: res["monthDay"],
		Month:    res["month"],
		WeekDay:  res["weekDay"],
	}

	err := checkSchedule(result)
	if err != nil {
		return Schedule{}, err
	}
	return result, nil
}
