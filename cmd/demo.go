package main

import (
	"fmt"

	cron "github.com/mlavergn/gocron"
)

func main() {
	fmt.Println("Go Cron Demo")

	pack := cron.NewCron()
	pack.Start()
}
