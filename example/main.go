package main

import (
	"fmt"

	"github.com/artem-shestakov/go-cron"
)

func main() {

	cron := cron.Cron("0/2 23/2 1/2 7/2 1/2")
	fmt.Println(cron.Minutes)
}
