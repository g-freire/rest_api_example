package common

import (
	"testing"
	"time"
)

type addTestValid struct {
	arg1, arg2 time.Time
	expected     bool
}

type addTestUpToDate struct {
	arg1 time.Time
	expected     bool
}

var addTests = []addTestValid{
	{time.Now(), time.Now().Add(time.Hour * 24), true},
	{time.Now(), time.Now().Add(time.Second * 1), true},
	{time.Now().Add(time.Hour * 24), time.Now(), false},
	{time.Now().Add(time.Second * 1), time.Now(), false},
}

var addTestsUpToDate = []addTestUpToDate{
	{time.Now().Add(time.Hour * 24), true},
	{time.Now().Add(time.Second * 1), true},
	{time.Now().Add(time.Hour * -100), false},
	{time.Now().Add(time.Second * -100), false},
}

func TestCheckTimestampIsValid(t *testing.T) {
	for _, test := range addTests {
		if output := CheckTimestampIsValid(test.arg1, test.arg2); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}
}

func TestCheckTimestampIsUpToDate(t *testing.T) {
	for _, test := range addTestsUpToDate {
		if output := CheckTimestampIsUpToDate(test.arg1); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}
}