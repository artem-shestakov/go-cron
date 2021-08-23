package cron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ()

type CronSchedule struct {
	Dows    map[int][]string
	Minutes Time
	Hours   Time
}

type Time struct {
	min      int
	max      int
	interval int
}

func Cron(cronSchedule string) *CronSchedule {
	schedule := strings.Split(cronSchedule, " ")
	dows := parseDow(schedule[4])
	return &CronSchedule{
		Dows: dows,
	}
}

func parseDow(dowSchedule string) map[int][]string {
	var separatorDows = regexp.MustCompile(`^[0-6](,([0-6]|[0-6]-[0-6]))*$`)
	dows := make(map[int][]string)
	switch {
	case dowSchedule == "*":
		for numOfWeek := 0; numOfWeek <= 6; numOfWeek++ {
			dows[numOfWeek] = make([]string, 0)
		}
	case separatorDows.MatchString(dowSchedule):
		for _, dow := range strings.Split(dowSchedule, ",") {
			if isRange(dow) {
				min, max := splitRange(dow)
				for numOfWeek := min; numOfWeek <= max; numOfWeek++ {
					fmt.Println(numOfWeek)
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
			fmt.Println(numOfWeek)
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
				fmt.Println(numOfWeek)
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

func isRange(rangeString string) bool {
	return regexp.MustCompile(`^[0-6]-[0-6]$`).MatchString(rangeString)
}

func isDow(dowString string) bool {
	return regexp.MustCompile(`^[0-6]$`).MatchString(dowString)
}

func splitRange(rangeString string) (int, int) {
	min, _ := strconv.Atoi(strings.Split(rangeString, "-")[0])
	max, _ := strconv.Atoi(strings.Split(rangeString, "-")[len(strings.Split(rangeString, "-"))-1])
	return min, max
}
