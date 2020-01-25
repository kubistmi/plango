package schedule

import (
	"fmt"
	"time"

	"github.com/kubistmi/plango/utils"
)

// Next ...
func (s Schedule) Next(after time.Time) (time.Time, error) {

	var reset time.Time
	next := after

	next, reset = s.NextTime(next)
	date, err := s.NextDate(next)
	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), fmt.Errorf("unable to find the date satisfying the schedule")
	}

	if date != next {
		return time.Date(date.Year(), date.Month(), date.Day(), reset.Hour(), reset.Minute(), reset.Second(), 0, time.Local), nil
	}
	return date, nil

}

// NextTime ...
func (s Schedule) NextTime(next time.Time) (time.Time, time.Time) {

	hours := findCandidates(s.Hour, next.Hour(), 24)
	minutes := findCandidates(s.Minute, next.Minute(), 60)
	seconds := findCandidates(s.Second, next.Second(), 60)
	minTime := time.Date(next.Year(), next.Month(), next.Day()+1, s.Hour.min(0), s.Minute.min(0), s.Second.min(0), 0, time.Local)

	nowSec := next.Hour()*3600 + next.Minute()*60 + next.Second()
	nxtHour, nxtMin, nxtSec, err := walkTime(hours, minutes, seconds, nowSec)
	if err != nil {
		return minTime, minTime
	}
	return time.Date(next.Year(), next.Month(), next.Day(), nxtHour, nxtMin, nxtSec, 0, time.Local), minTime
}

// NextDate ...
func (s Schedule) NextDate(next time.Time) (time.Time, error) {
	// TODO: should this be a config variable?
	iter := 50

	switch s.WeekDay.(type) {

	case partAny:
		// the easy part, for non-specific weekDay, just go through the calendar
		var nxtMday, nxtMonth int
		var mdShift, mShift int

		for i := 0; i < iter; i++ {
			nxtMday, mdShift = s.MonthDay.compareTime(next.Day())
			next = next.AddDate(0, mdShift, 0)

			nxtMonth, mShift = s.Month.compareTime(int(next.Month()))
			next = next.AddDate(mShift, nxtMonth-int(next.Month()), nxtMday-next.Day())

			if int(next.Month()) != nxtMonth {
				next = next.AddDate(0, 0, 1)
				continue
			} else {
				return next, nil
			}
		}

	default:

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
				if s.MonthDay.isin(next.Day()) && s.Month.isin(int(next.Month())) {
					return next, nil
				}
				next = next.AddDate(0, 0, 1)
				continue
			} else if wdMshift == 1 {
				// if Schedule.Month too high, then jump to the next year and ...
				// else to the max(Schedule.Month, 1) and ...
				// ... lowest of the Schedule.MonthDay
				next = next.AddDate(1, s.Month.min(1)-int(next.Month()), s.MonthDay.min(1)-next.Day())
			} else {
				monthShift := []int{wdMonth - int(next.Month()), wdShift}
				next = next.AddDate(0, utils.FindMax(monthShift), s.MonthDay.min(1)-next.Day())
			}
		}
	}
	//! if the date cannot be found, not sure this is reachable
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), fmt.Errorf("unable to find the date satisfying the schedule")
}

func findCandidates(p part, next int, check int) []int {
	candidates := make([]int, 0, 60)
	switch v := p.(type) {
	case partList:
		for _, val := range v.List {
			if val >= next {
				candidates = append(candidates, val)
			}
		}
		if len(candidates) >= 3 {
			candidates = candidates[:3]
		} else {
			candidates = append([]int{v.min(0)}, candidates...)
		}
	case partAny:
		for i := 0; i < 4; i++ {
			candidates = append(candidates, (next+i)%check)
		}
	}
	return candidates
}

func walkTime(hours, minutes, seconds []int, nowSec int) (int, int, int, error) {
	for _, h := range hours {
		for _, m := range minutes {
			for _, s := range seconds {
				if h*3600+m*60+s >= nowSec {
					return h, m, s, nil
				}
			}
		}
	}
	return 0, 0, 0, fmt.Errorf("Shift the day")
}
