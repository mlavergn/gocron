package main

import (
	"fmt"
	"log"
	"time"

	cron "github.com/mlavergn/gocron"
)

func main() {
	fmt.Println("Go Cron Demo")

	pack := cron.NewCron()
	// job := cron.NewJob(1, -1, -1, -1, -1, "demo", "/bin/sh", "-c", "date >> /tmp/gocron.log")
	job := cron.NewJobFn(1, -1, -1, -1, -1, "demo", func() {
		log.Println("Hello")
	})
	job.MinuteEvery = true
	pack.Add(job)
	pack.List()
	fmt.Println(job.Next(), time.Until(job.Next()))
	pack.Start()
	// arbtrary duration of 1 hour
	<-time.After(1 * time.Hour)
}
