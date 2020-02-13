package cron

import (
	"fmt"
	"testing"
	"time"
)

func TestNextTime(t *testing.T) {
	job := NewJob(5, 4, -1, -1, -1, "ls")
	println(job.String())
	fmt.Println(time.Now(), job.NextTime())

	job = NewJob(1, -1, -1, -1, -1, "ls")
	println(job.String())
	fmt.Println(time.Now(), job.NextTime())

	job = NewJob(15, -1, -1, -1, -1, "ls")
	println(job.String())
	fmt.Println(time.Now(), job.NextTime())

	job = NewJob(55, -1, -1, -1, -1, "ls")
	println(job.String())
	fmt.Println(time.Now(), job.NextTime())

	job = NewJob(-1, 7, -1, -1, -1, "ls")
	println(job.String())
	fmt.Println(time.Now(), job.NextTime())

	job = NewJob(-1, 22, -1, -1, -1, "ls")
	println(job.String())
	fmt.Println(time.Now(), job.NextTime())

	job = NewJob(0, 2, -1, -1, 2, "ls")
	job.DOWEvery = true
	println(job.String())
	fmt.Println(time.Now(), job.NextTime())
}
