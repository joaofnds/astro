package habit

import (
	"astro/date"
	"testing"
)

func TestMomentum(t *testing.T) {
	activities := []Activity{
		{CreatedAt: date.Today().AddDate(0, 0, -9)}, // +1
		//                                     skipped: -1 = 0
		{CreatedAt: date.Today().AddDate(0, 0, -7)}, // +1 = 1
		//                                     skipped: -1 = 0
		//                                     skipped: -1 = 0  <- min is 0
		{CreatedAt: date.Today().AddDate(0, 0, -4)}, // +1 = 1
		{CreatedAt: date.Today().AddDate(0, 0, -4)}, // +1 = 2
		{CreatedAt: date.Today().AddDate(0, 0, -3)}, // +1 = 3
		//                                     skipped: -1 = 2
		{CreatedAt: date.Today().AddDate(0, 0, -1)}, // +1 = 3
	}

	expected := 3
	got := Momentum(activities)

	if got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func TestDigest(t *testing.T) {
	activities := []Activity{
		{CreatedAt: date.Today().AddDate(0, 0, -5)},
		{CreatedAt: date.Today().AddDate(0, 0, -4)},
		{CreatedAt: date.Today().AddDate(0, 0, -2)},
		{CreatedAt: date.Today().AddDate(0, 0, -1)},
	}
	got := Digest("run", activities)
	want := "run - streak: 2 days, momentum: 3"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
