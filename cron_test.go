package cron

import (
	"fmt"
	"testing"
	"time"
)

// Can be validated online via
// https://crontab.guru
func TestNextTime(t *testing.T) {
	verbose := true

	year := 2007
	month := time.Month(1)
	day := 9
	// dow := 2
	hour := 9
	min := 41
	loc := time.Now().Location()

	// the expected values need a known baseline to work
	// consistently, so this may create false negatives,
	// but now time can be useful for debugging
	useNow := false
	if useNow {
		nowTime := time.Now()
		year = nowTime.Year()
		month = nowTime.Month()
		day = nowTime.Day()
		hour = nowTime.Hour()
		min = nowTime.Minute()
	}

	testTime := time.Date(year, month, day, hour, min, 0, 0, loc)
	debugTime = &testTime

	if verbose {
		fmt.Println("Base time:", debugTime)
	}

	// 5 4 * * *
	job := NewJob(5, 4, -1, -1, -1, "test", "ls")
	actualString := job.TimeString()
	expectedString := "5 4 * * *"
	if actualString != expectedString {
		t.Fatal("TestNextTime unexpected String result", actualString, expectedString)
	}
	actual := job.Next()
	expected := time.Date(year, month, day+1, 4, 5, 0, 0, loc)
	if actual != expected {
		t.Fatal("TestNextTime unexpected Time result", actual, expected)
	}
	if verbose {
		fmt.Println(expectedString, "\t", actual)
	}

	// 1 * * * *
	job = NewJob(1, -1, -1, -1, -1, "test", "ls")
	actualString = job.TimeString()
	expectedString = "1 * * * *"
	if actualString != expectedString {
		t.Fatal("TestNextTime unexpected String result", actualString, expectedString)
	}
	actual = job.Next()
	expected = time.Date(year, month, day, hour+1, 1, 0, 0, loc)
	if actual != expected {
		t.Fatal("TestNextTime unexpected Time result", expectedString, actual, expected)
	}
	if verbose {
		fmt.Println(expectedString, "\t", actual)
	}

	// 15 * * * *
	job = NewJob(15, -1, -1, -1, -1, "test", "ls")
	actualString = job.TimeString()
	expectedString = "15 * * * *"
	if actualString != expectedString {
		t.Fatal("TestNextTime unexpected String result", actualString, expectedString)
	}
	actual = job.Next()
	expected = time.Date(year, month, day, hour+1, 15, 0, 0, loc)
	if actual != expected {
		t.Fatal("TestNextTime unexpected Time result", expectedString, actual, expected)
	}
	if verbose {
		fmt.Println(expectedString, "\t", actual)
	}

	// 55 * * * *
	job = NewJob(55, -1, -1, -1, -1, "test", "ls")
	actualString = job.TimeString()
	expectedString = "55 * * * *"
	if actualString != expectedString {
		t.Fatal("TestNextTime unexpected String result", actualString, expectedString)
	}
	actual = job.Next()
	expected = time.Date(year, month, day, hour, 55, 0, 0, loc)
	if actual != expected {
		t.Fatal("TestNextTime unexpected Time result", expectedString, actual, expected)
	}
	if verbose {
		fmt.Println(expectedString, "\t", actual)
	}

	// * 7 * * *
	job = NewJob(-1, 7, -1, -1, -1, "test", "ls")
	actualString = job.TimeString()
	expectedString = "* 7 * * *"
	if actualString != expectedString {
		t.Fatal("TestNextTime unexpected String result", actualString, expectedString)
	}
	actual = job.Next()
	expected = time.Date(year, month, day+1, 7, 0, 0, 0, loc)
	if actual != expected {
		t.Fatal("TestNextTime unexpected Time result", expectedString, actual, expected)
	}
	if verbose {
		fmt.Println(expectedString, "\t", actual)
	}

	// * 22 * * *
	job = NewJob(-1, 22, -1, -1, -1, "test", "ls")
	actualString = job.TimeString()
	expectedString = "* 22 * * *"
	if actualString != expectedString {
		t.Fatal("TestNextTime unexpected String result", actualString, expectedString)
	}
	actual = job.Next()
	expected = time.Date(year, month, day, 22, 0, 0, 0, loc)
	if actual != expected {
		t.Fatal("TestNextTime unexpected Time result", expectedString, actual, expected)
	}
	if verbose {
		fmt.Println(expectedString, "\t", actual)
	}

	// 0 2 * * */3
	job = NewJob(0, 2, -1, -1, 3, "test", "ls")
	job.DOWEvery = true
	actualString = job.TimeString()
	expectedString = "0 2 * * */3"
	if actualString != expectedString {
		t.Fatal("TestNextTime unexpected String result", actualString, expectedString)
	}
	actual = job.Next()
	expected = time.Date(year, month, day+2, 2, 0, 0, 0, loc)
	if actual != expected {
		t.Fatal("TestNextTime unexpected Time result", expectedString, actual, expected)
	}
	if verbose {
		fmt.Println(expectedString, "\t", actual)
	}
}
