package cron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type CronSchedule struct {
	Day      map[int][]string
	Month    map[int][]string
	Dows     map[int][]string
	Schedule map[int]map[int][]string
}

type Unit struct {
	Min            int
	Max            int
	SeparatorReqex string
	IntervalReqex  string
}

var dowUnit = Unit{
	Min:            1,
	Max:            7,
	SeparatorReqex: `^[0-6](,([0-6]|[0-6]-[0-6]))*$`,
	IntervalReqex:  `^[0-6]-[0-6]$`,
}

var monthUnit = Unit{
	Min:            1,
	Max:            12,
	SeparatorReqex: `^(0?[1-9]|1[012])(,(0?[1-9]|1[012]))*$`,
	IntervalReqex:  `^(0?[1-9]|1[012])-(0?[1-9]|1[012])$`,
}

var daysOfMonths = map[int]int{
	1: 31, 2: 29, 3: 31, 4: 30,
	5: 31, 6: 30, 7: 31, 8: 31,
	9: 30, 10: 31, 11: 30, 12: 31,
}

func Cron(cronSchedule string) *CronSchedule {
	schedule := strings.Split(cronSchedule, " ")
	day := parseDays(schedule[2])
	month := parseMonth(schedule[3])
	dows := parseDow(schedule[4])
	schedule2 := scheduling(month, schedule[2])
	return &CronSchedule{
		Day:      day,
		Month:    month,
		Dows:     dows,
		Schedule: schedule2,
	}
}

func scheduling(months map[int][]string, daySchedule string) map[int]map[int][]string {
	var total = make(map[int]map[int][]string)
	for month, _ := range months {
		total[month] = parseDays(daySchedule)
	}
	return total
}

func parseDow(dowSchedule string) map[int][]string {
	var separatorDows = regexp.MustCompile(`^[0-6](,([0-6]|[0-6]-[0-6]))*$`)
	dows := make(map[int][]string)
	switch {
	case dowSchedule == "*":
		for numOfWeek := 0; numOfWeek < 7; numOfWeek++ {
			dows[numOfWeek] = make([]string, 0)
		}
	case separatorDows.MatchString(dowSchedule):
		for _, dow := range strings.Split(dowSchedule, ",") {
			if isRange(dow) {
				min, max := splitRange(dow)
				for numOfWeek := min; numOfWeek <= max; numOfWeek++ {
					dows[numOfWeek] = make([]string, 0)
				}
				continue
			}
			numOfWeek, _ := strconv.Atoi(dow)
			dows[numOfWeek] = make([]string, 0)
		}
	case isRange(dowSchedule):
		min, max := splitRange(dowSchedule)
		for numOfWeek := min; numOfWeek <= max; numOfWeek++ {
			dows[numOfWeek] = make([]string, 0)
		}
	case len(strings.Split(dowSchedule, "/")) == 2:
		value := strings.Split(dowSchedule, "/")[0]
		interval, _ := strconv.Atoi(strings.Split(dowSchedule, "/")[1])
		if value == "*" {
			for numOfWeek := 0; numOfWeek < 7; numOfWeek += interval {
				dows[numOfWeek] = make([]string, 0)
			}
		} else if isRange(value) {
			min, max := splitRange(value)
			fmt.Println(min, max)
			for numOfWeek := min; numOfWeek <= max; numOfWeek += interval {
				dows[numOfWeek] = make([]string, 0)
			}
		} else if isDow(value) {
			startDow, _ := strconv.Atoi(value)
			for numOfWeek := startDow; numOfWeek <= 7; numOfWeek += interval {
				if numOfWeek == 7 {
					dows[0] = make([]string, 0)
					continue
				}
				dows[numOfWeek] = make([]string, 0)
			}
		}
	}
	return dows
}

func parseMonth(monthSchedule string) map[int][]string {
	var separatorMonths = regexp.MustCompile(`^(0?[1-9]|1[012])(,(0?[1-9]|1[012]))*$`)
	months := make(map[int][]string)
	switch {
	case monthSchedule == "*":
		for numOfMonth := 1; numOfMonth <= 12; numOfMonth++ {
			months[numOfMonth] = make([]string, 0)
		}
	case separatorMonths.MatchString(monthSchedule):
		for _, month := range strings.Split(monthSchedule, ",") {
			if isRange(month) {
				min, max := splitRange(month)
				for numOfMonth := min; numOfMonth <= max; numOfMonth++ {
					months[numOfMonth] = make([]string, 0)
				}
				continue
			}
			numOfMonth, _ := strconv.Atoi(month)
			months[numOfMonth] = make([]string, 0)
		}
	case regexp.MustCompile(`^(0?[1-9]|1[012])-(0?[1-9]|1[012])$`).MatchString(monthSchedule):
		min, max := splitRange(monthSchedule)
		for numOfMonth := min; numOfMonth <= max; numOfMonth++ {
			months[numOfMonth] = make([]string, 0)
		}
	case len(strings.Split(monthSchedule, "/")) == 2:
		value := strings.Split(monthSchedule, "/")[0]
		interval, _ := strconv.Atoi(strings.Split(monthSchedule, "/")[1])
		if value == "*" {
			for numOfMonth := 1; numOfMonth <= 12; numOfMonth += interval {
				months[numOfMonth] = make([]string, 0)
			}
		} else if regexp.MustCompile(`^(0?[1-9]|1[012])-(0?[1-9]|1[012])$`).MatchString(value) {
			min, max := splitRange(value)
			fmt.Println(min, max)
			fmt.Println(min, max)
			for numOfMonth := min; numOfMonth <= max; numOfMonth += interval {
				months[numOfMonth] = make([]string, 0)
			}
		} else if isMonth(value) {
			startMonth, _ := strconv.Atoi(value)
			for numOfMonth := startMonth; numOfMonth <= 12; numOfMonth += interval {
				months[numOfMonth] = make([]string, 0)
			}
		}
	}
	return months
}

func parseDays(daySchedule string) map[int][]string {
	var separatorMonths = regexp.MustCompile(`^(0?[1-9]|1[0-9]|2[0-9]|3[01])(,(0?[1-9]|1[0-9]|2[0-9]|3[01]))*$`)
	days := make(map[int][]string)
	switch {
	case daySchedule == "*":
		for dayOfMonth := 1; dayOfMonth <= 31; dayOfMonth++ {
			days[dayOfMonth] = make([]string, 0)
		}
	case separatorMonths.MatchString(daySchedule):
		for _, month := range strings.Split(daySchedule, ",") {
			if isRange(month) {
				min, max := splitRange(month)
				for dayOfMonth := min; dayOfMonth <= max; dayOfMonth++ {
					days[dayOfMonth] = make([]string, 0)
				}
				continue
			}
			dayOfMonth, _ := strconv.Atoi(month)
			days[dayOfMonth] = make([]string, 0)
		}
	case regexp.MustCompile(`^(0?[1-9]|1[0-9]|2[0-9]|3[01])-(0?[1-9]|1[0-9]|2[0-9]|3[01])$`).MatchString(daySchedule):
		min, max := splitRange(daySchedule)
		for dayOfMonth := min; dayOfMonth <= max; dayOfMonth++ {
			days[dayOfMonth] = make([]string, 0)
		}
	case len(strings.Split(daySchedule, "/")) == 2:
		value := strings.Split(daySchedule, "/")[0]
		interval, _ := strconv.Atoi(strings.Split(daySchedule, "/")[1])
		if value == "*" {
			for dayOfMonth := 1; dayOfMonth <= 31; dayOfMonth += interval {
				days[dayOfMonth] = make([]string, 0)
			}
		} else if regexp.MustCompile(`^(0?[1-9]|1[0-9]|2[0-9]|3[01])-(0?[1-9]|1[0-9]|2[0-9]|3[01])$`).MatchString(value) {
			min, max := splitRange(value)
			fmt.Println(min, max)
			fmt.Println(min, max)
			for dayOfMonth := min; dayOfMonth <= max; dayOfMonth += interval {
				days[dayOfMonth] = make([]string, 0)
			}
		} else if isDay(value) {
			startDay, _ := strconv.Atoi(value)
			for dayOfMonth := startDay; dayOfMonth <= 31; dayOfMonth += interval {
				days[dayOfMonth] = make([]string, 0)
			}
		}
	}
	return days
}

// func parse(monthSchedule string, unit Unit) map[int][]string {
// 	var separatorMonths = regexp.MustCompile(unit.SeparatorReqex)
// 	months := make(map[int][]string)
// 	switch {
// 	case monthSchedule == "*":
// 		for numOfMonth := unit.Min; numOfMonth <= unit.Max; numOfMonth++ {
// 			months[numOfMonth] = make([]string, 0)
// 		}
// 	case separatorMonths.MatchString(monthSchedule):
// 		for _, month := range strings.Split(monthSchedule, ",") {
// 			if isRange(month) {
// 				min, max := splitRange(month)
// 				for numOfMonth := min; numOfMonth <= max; numOfMonth++ {
// 					months[numOfMonth] = make([]string, 0)
// 				}
// 				continue
// 			}
// 			numOfMonth, _ := strconv.Atoi(month)
// 			months[numOfMonth] = make([]string, 0)
// 		}
// 	case regexp.MustCompile(unit.IntervalReqex).MatchString(monthSchedule):
// 		min, max := splitRange(monthSchedule)
// 		for numOfMonth := min; numOfMonth <= max; numOfMonth++ {
// 			months[numOfMonth] = make([]string, 0)
// 		}
// 	case len(strings.Split(monthSchedule, "/")) == 2:
// 		value := strings.Split(monthSchedule, "/")[0]
// 		interval, _ := strconv.Atoi(strings.Split(monthSchedule, "/")[1])
// 		if value == "*" {
// 			for numOfMonth := unit.Min; numOfMonth <= unit.Max; numOfMonth += interval {
// 				months[numOfMonth] = make([]string, 0)
// 			}
// 		} else if regexp.MustCompile(unit.IntervalReqex).MatchString(value) {
// 			min, max := splitRange(value)
// 			fmt.Println(min, max)
// 			fmt.Println(min, max)
// 			for numOfMonth := min; numOfMonth <= max; numOfMonth += interval {
// 				months[numOfMonth] = make([]string, 0)
// 			}
// 		} else if isMonth(value) {
// 			startMonth, _ := strconv.Atoi(value)
// 			for numOfMonth := startMonth; numOfMonth <= unit.Max; numOfMonth += interval {
// 				months[numOfMonth] = make([]string, 0)
// 			}
// 		}
// 	}
// 	return months
// }

func isRange(rangeString string) bool {
	return regexp.MustCompile(`^[0-6]-[0-6]$`).MatchString(rangeString)
}

func isDow(dowString string) bool {
	return regexp.MustCompile(`^[0-6]$`).MatchString(dowString)
}

func isMonth(dowString string) bool {
	return regexp.MustCompile(`^([1-9]|1[012])$`).MatchString(dowString)
}

func isDay(dowString string) bool {
	return regexp.MustCompile(`^(0?[1-9]|1[0-9]|2[0-9]|3[01])$`).MatchString(dowString)
}

func splitRange(rangeString string) (int, int) {
	min, _ := strconv.Atoi(strings.Split(rangeString, "-")[0])
	max, _ := strconv.Atoi(strings.Split(rangeString, "-")[len(strings.Split(rangeString, "-"))-1])
	return min, max
}
