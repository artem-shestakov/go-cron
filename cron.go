package cron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Schedule struct {
	Minutes []int
	Hours   []int
	Days    []int
	Months  []int
	Dows    []int
}

type Unit struct {
	min     int
	max     int
	pattern string
}

var minuteUnit = Unit{
	min:     0,
	max:     59,
	pattern: `(0?[1-9]|[1-5][0-9])`,
}

var hourUnit = Unit{
	min:     0,
	max:     23,
	pattern: `(0?[1-9]|1[0-9]|2[0-3])`,
}

var dowUnit = Unit{
	min:     0,
	max:     6,
	pattern: `[0-6]`,
}

var monthUnit = Unit{
	min:     1,
	max:     12,
	pattern: `(0?[1-9]|1[012])`,
}

var dayUnit = Unit{
	min:     1,
	max:     31,
	pattern: `(0?[1-9]|[12][0-9]|3[01])`,
}

var daysOfMonths = map[int]int{
	1: 31, 2: 29, 3: 31, 4: 30,
	5: 31, 6: 30, 7: 31, 8: 31,
	9: 30, 10: 31, 11: 30, 12: 31,
}

func Cron(cronSchedule string) *Schedule {
	schedule := strings.Split(cronSchedule, " ")
	return &Schedule{
		Minutes: parse(minuteUnit, schedule[0]),
		Hours:   parse(hourUnit, schedule[1]),
		Days:    parse(dayUnit, schedule[2]),
		Months:  parse(monthUnit, schedule[3]),
		Dows:    parse(dowUnit, schedule[4]),
	}
}

func parse(unitInfo Unit, unitSchedule string) []int {
	units := make([]int, 0)
	switch {
	// '*'
	case unitSchedule == "*":
		for numOfUnit := unitInfo.min; numOfUnit <= unitInfo.max; numOfUnit++ {
			units = append(units, numOfUnit)
		}
	// '1,2,3' or '1'
	case regexp.MustCompile(fmt.Sprintf("^%s(,%s)*$", unitInfo.pattern, unitInfo.pattern)).MatchString(unitSchedule):
		for _, unit := range strings.Split(unitSchedule, ",") {
			if isRange(unitInfo.pattern, unit) {
				min, max := splitRange(unit)
				for numOfUnit := min; numOfUnit <= max; numOfUnit++ {
					units = append(units, numOfUnit)
				}
				continue
			}
			numOfUnit, _ := strconv.Atoi(unit)
			units = append(units, numOfUnit)
		}
	// '1-5'
	case isRange(unitInfo.pattern, unitSchedule):
		min, max := splitRange(unitSchedule)
		for numOfUnit := min; numOfUnit <= max; numOfUnit++ {
			units = append(units, numOfUnit)
		}
	// '*/2' or '4/2' or '1-5/2'
	case len(strings.Split(unitSchedule, "/")) == 2:
		value := strings.Split(unitSchedule, "/")[0]
		interval, _ := strconv.Atoi(strings.Split(unitSchedule, "/")[1])
		if value == "*" {
			startUnit := unitInfo.min
			if startUnit == 0 {
				startUnit = 1
			}
			for numOfUnit := startUnit; numOfUnit <= unitInfo.max; numOfUnit += interval {
				units = append(units, numOfUnit)
			}
		} else if value == "0" {
			units = append(units, 0)
		} else if isRange(unitInfo.pattern, value) {
			min, max := splitRange(value)
			for numOfUnit := min; numOfUnit <= max; numOfUnit += interval {
				units = append(units, numOfUnit)
			}
		} else if isUnit(unitInfo.pattern, value) {
			startUnit, _ := strconv.Atoi(value)
			endUnit := unitInfo.max
			if startUnit == 0 {
				startUnit = 1
				endUnit++
			} else if startUnit != 0 && unitInfo.min == 0 {
				endUnit++
			}
			for numOfUnit := startUnit; numOfUnit <= endUnit; numOfUnit += interval {
				if numOfUnit == unitInfo.max+1 {
					numOfUnit = 0
					units = append(units, numOfUnit)
					break
				}
				units = append(units, numOfUnit)
			}
		}
	}
	return units
}

func isRange(pattern, expression string) bool {
	return regexp.MustCompile(fmt.Sprintf(`^%s-%s$`, pattern, pattern)).MatchString(expression)
}

func isUnit(pattern, expression string) bool {
	return regexp.MustCompile(fmt.Sprintf(`^%s$`, pattern)).MatchString(expression)
}

func splitRange(rangeString string) (int, int) {
	min, _ := strconv.Atoi(strings.Split(rangeString, "-")[0])
	max, _ := strconv.Atoi(strings.Split(rangeString, "-")[len(strings.Split(rangeString, "-"))-1])
	return min, max
}
