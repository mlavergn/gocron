package cron

import (
	"log"
	oslog "log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Version export
const Version = "0.1.0"

// logger stand-in
var dlog *oslog.Logger

// DEBUG toggle
var DEBUG = false

// Job export
type Job struct {
	Minute      int // 0-59
	MinuteEvery bool
	Hour        int // 0-23
	HourEvery   bool
	DOM         int // 1-31
	DOMEvery    bool
	Month       int // 1-12
	MonthEvery  bool
	DOW         int // 0-7 - sunday loops
	DOWEvery    bool
	Command     string
}

// NewJob init
// -1 == *
func NewJob(min int, hour int, dom int, month int, dow int, command string) *Job {
	id := &Job{
		Minute:  min,
		Hour:    hour,
		DOM:     dom,
		Month:   month,
		DOW:     dow,
		Command: command,
	}

	if min > 59 || hour > 23 || dom > 31 || month > 12 || dow > 7 {
		return nil
	}

	return id
}

// String export
func (id *Job) String() string {
	fmtSegment := func(value int, every bool) string {
		prefix := ""
		if every {
			prefix = "*/"
		}
		segment := "*"
		if value != -1 {
			segment = strconv.Itoa(value)
		}
		return prefix + segment + " "
	}
	var result strings.Builder
	result.WriteString(fmtSegment(id.Minute, id.MinuteEvery))
	result.WriteString(fmtSegment(id.Hour, id.HourEvery))
	result.WriteString(fmtSegment(id.DOM, id.DOMEvery))
	result.WriteString(fmtSegment(id.Month, id.MonthEvery))
	result.WriteString(fmtSegment(id.DOW, id.DOWEvery))

	return result.String()
}

// NextTime export
// Implements the cron time interval logic
// https://crontab.guru
func (id *Job) NextTime() time.Time {
	now := time.Now()

	// capture the current values
	year, month, day := now.Date()
	hour := now.Hour()
	minute := now.Minute()

	// cron does not track seconds or lower
	second := 0
	nanosecond := 0

	// override values

	// minutes
	if id.Minute != -1 {
		if id.Minute < minute && !id.MinuteEvery {
			// overflow to next hour
			hour++
		}
		if id.MinuteEvery {
			minute += id.Minute
		} else {
			minute = id.Minute
		}
	} else {
		minute++
	}

	// hour
	if id.Hour != -1 {
		// if we're in the wrong hour, reset min to 0 if *
		if id.Hour != hour && id.Minute == -1 {
			minute = 0
		}
		if id.Hour < hour && !id.HourEvery {
			day++
		}
		if id.HourEvery {
			hour += id.Hour
		} else {
			hour = id.Hour
		}
	}

	if id.Month != -1 {
		if int(id.Month) < int(month) && !id.MonthEvery {
			year++
		}
		if id.MonthEvery {
			month = time.Month(int(month) + id.Month)
		} else {
			month = time.Month(id.Month)
		}
	}

	if id.DOM != -1 {
		if id.DOM < day && !id.DOMEvery {
			month = time.Month(id.Month + 1)
		}
		if id.DOMEvery {
			day += id.DOM
		} else {
			day = id.DOM
		}
	}

	// calculate the next date
	next := time.Date(year, month, day,
		hour, minute, second, nanosecond, now.Location())

	// adjust for DOW
	if id.DOW != -1 {
		if !id.DOMEvery {
			dow := int(next.Weekday()) - id.DOW
			if dow < 0 {
				dow += 7
			}
			next.Add(time.Duration(dow) * 24 * time.Hour)
		} else {
			next.Add(time.Duration(id.DOW) * 24 * time.Hour)
		}
	}

	return next
}

// Cron export
type Cron struct {
	Jobs []Job
}

// NewCron init
func NewCron() *Cron {
	return &Cron{}
}

// Add export
func (id *Cron) Add(job *Job) {
	id.Jobs = append(id.Jobs, *job)
}

func (id *Cron) schedule(job Job) {
	duration := time.Until(job.NextTime())
	fn := func() {
		cmd := exec.Command(job.Command)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
		// requeue the job for it's next run (soonest 1 min)
		<-time.After(1 * time.Second)
		id.schedule(job)
	}
	time.AfterFunc(duration, fn)
}

// Start export
func (id *Cron) Start() {
	for _, job := range id.Jobs {
		id.schedule(job)
	}
}
