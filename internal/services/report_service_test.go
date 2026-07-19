package services

import (
	"testing"
	"time"
)

func TestReportService_MonthNames(t *testing.T) {
	monthNamesIndo := map[time.Month]string{
		time.January:   "Jan",
		time.February:  "Feb",
		time.March:     "Mar",
		time.April:     "Apr",
		time.May:       "Mei",
		time.June:      "Jun",
		time.July:      "Jul",
		time.August:    "Agu",
		time.September: "Sep",
		time.October:   "Okt",
		time.November:  "Nov",
		time.December:  "Des",
	}

	if monthNamesIndo[time.May] != "Mei" {
		t.Errorf("expected Mei, got %s", monthNamesIndo[time.May])
	}
	if monthNamesIndo[time.August] != "Agu" {
		t.Errorf("expected Agu, got %s", monthNamesIndo[time.August])
	}
}
