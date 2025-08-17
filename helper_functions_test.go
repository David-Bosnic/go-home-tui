package main

import (
	"fmt"
	"testing"
	"time"
)

func TestEventTimeFormatter(t *testing.T) {
	var tests = []struct {
		stringTime    string
		formattedTime time.Time
	}{
		{"hello", time.Now()},
	}
	ans := EventTimeFormatter("beeb")
	if ans != time.Time.Day(time.Now()) {
		t.Error("Did not work :(")
	}
}
