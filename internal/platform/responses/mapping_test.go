package responses

import (
	"strconv"
	"testing"
	"time"
)

func TestMapSlice(t *testing.T) {
	got := MapSlice([]int(nil), strconv.Itoa)
	if got == nil || len(got) != 0 {
		t.Fatalf("nil input: want empty non-nil slice, got %#v", got)
	}
	got = MapSlice([]int{1, 2, 3}, strconv.Itoa)
	if len(got) != 3 || got[0] != "1" || got[2] != "3" {
		t.Fatalf("unexpected result %#v", got)
	}
}

func TestFormatTime(t *testing.T) {
	ts := time.Date(2026, 11, 16, 10, 30, 0, 0, time.UTC)
	if got := FormatTime(ts); got != "2026-11-16T10:30:00Z" {
		t.Fatalf("got %q", got)
	}
	if FormatTimePtr(nil) != nil {
		t.Fatal("nil time should give nil")
	}
	if p := FormatTimePtr(&ts); p == nil || *p != "2026-11-16T10:30:00Z" {
		t.Fatalf("got %v", p)
	}
}
