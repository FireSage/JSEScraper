package main

import "testing"

func TestGetDollarValueAsInt(t *testing.T) {
	got := getDollarValueAsInt("$123.45")
	want := 12345

	if got != int64(want) {
		t.Errorf("got %q want %q", got, want)
	}
}
