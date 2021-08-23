package main

import (
	"fmt"

	"github.com/artem-shestakov/go-cron"
)

func main() {
	cron := cron.Cron("* * * 5/2 *")
	fmt.Println(cron.Schedule)
}
