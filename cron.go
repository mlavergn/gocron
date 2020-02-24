package cron

import (
	"io/ioutil"
	"log"
	oslog "log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Version export
const Version = "0.2.3"

// logger stand-in
var dlog *oslog.Logger

// DEBUG toggle
var DEBUG = false

var debugTime *time.Time

func init() {
	if DEBUG {
		dlog = oslog.New(os.Stderr, "", 0)
	} else {
		dlog = oslog.New(ioutil.Discard, "", 0)
	}
}

// -----------------------------------------------------------------------------
// Job

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
	Name        string
	Command     string
	Args        []string
	Fn          func()
}

// NewJob init
// -1 == *
func NewJob(min int, hour int, dom int, month int, dow int, name string, command string, args ...string) *Job {
	id := &Job{
		Minute:  min,
		Hour:    hour,
		DOM:     dom,
		Month:   month,
		DOW:     dow,
		Name:    name,
		Command: command,
		Args:    args,
	}

	if min > 59 || hour > 23 || dom > 31 || month > 12 || dow > 7 {
		return nil
	}

	return id
}

// NewJobFn init
// -1 == *
func NewJobFn(min int, hour int, dom int, month int, dow int, name string, fn func()) *Job {
	id := &Job{
		Minute: min,
		Hour:   hour,
		DOM:    dom,
		Month:  month,
		DOW:    dow,
		Name:   name,
		Fn:     fn,
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

	return strings.TrimSpace(result.String())
}

// Next export
// Implements the cron time interval logic
// BUG: Interval will be incorrect over leap events
func (id *Job) Next() time.Time {
	baseTime := time.Now()
	if debugTime != nil {
		baseTime = *debugTime
	}

	// capture the current values
	year, month, day := baseTime.Date()
	hour := baseTime.Hour()
	minute := baseTime.Minute()

	// cron does not track anything less than minutes
	second := 0
	nanosecond := 0

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

	// month
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

	// day of month
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

	// calculate the next time
	nextTime := time.Date(year, month, day,
		hour, minute, second, nanosecond, baseTime.Location())

	// day of week
	if id.DOW != -1 {
		if !id.DOWEvery {
			dow := int(nextTime.Weekday()) - id.DOW
			if dow < 0 {
				dow += 7
			}
			nextTime = nextTime.Add(time.Duration(dow) * 24 * time.Hour)
		} else {
			dow := id.DOW - 1
			// undo any added days
			if id.Hour > baseTime.Hour() {
				dow--
			}
			nextTime = nextTime.Add((time.Duration(dow) - 1) * 24 * time.Hour)
		}
	}

	return nextTime
}

// -----------------------------------------------------------------------------
// Cron

// Cron export
type Cron struct {
	jobs []*Job
}

// NewCron init
func NewCron() *Cron {
	return &Cron{
		jobs: []*Job{},
	}
}

// Add export
func (id *Cron) Add(job *Job) {
	dlog.Println("Cron.Add")
	id.jobs = append(id.jobs, job)
}

// Tab export
func (id *Cron) Tab() map[string]map[string]string {
	result := map[string]map[string]string{}
	for _, job := range id.jobs {
		result[job.Name] = map[string]string{
			"job":  job.String(),
			"next": job.Next().String(),
		}
	}
	return result
}

// List export
func (id *Cron) List() {
	for _, job := range id.jobs {
		log.Println(job.Name, job.String(), job.Next())
	}
}

func (id *Cron) schedule(job *Job) {
	dlog.Println("Cron.schedule")
	duration := time.Until(job.Next())
	fn := func() {
		if job.Fn != nil {
			log.Println("Running", job.Name, "as Fn")
			job.Fn()
		} else {
			log.Println("Running", job.Name, "as Command")
			cmd := exec.Command(job.Command, job.Args...)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}
		}
		// requeue the job for it's next run (soonest 1 min)
		<-time.After(1 * time.Second)
		id.schedule(job)
	}
	time.AfterFunc(duration, fn)
}

// Start export
func (id *Cron) Start() {
	dlog.Println("Cron.Start")
	for _, job := range id.jobs {
		id.schedule(job)
	}
}
