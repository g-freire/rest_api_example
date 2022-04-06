package class

import (
	"testing"
	"time"
)

type addTest struct {
	arg1, arg2 time.Time
	expected     bool
}

var addTests = []addTest{
	addTest{time.Now(), time.Now().Add(time.Hour * 24), true},
	addTest{time.Now(), time.Now().Add(time.Second * 1), true},
	addTest{time.Now().Add(time.Hour * 24), time.Now(), false},
	addTest{time.Now().Add(time.Second * 1), time.Now(), false},
}

func TestCheckTimestampIsValid(t *testing.T) {
	for _, test := range addTests {
		if output := CheckTimestampIsValid(test.arg1, test.arg2); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}
}