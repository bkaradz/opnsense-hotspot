package printing

import (
	"fmt"
	"time"
)

func FormatDate(t time.Time) string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	dayName := days[int(t.Weekday())]
	day := t.Day()
	month := months[int(t.Month())-1]
	year := t.Year()

	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	period := "AM"
	if hour >= 12 {
		period = "PM"
	}
	hour12 := hour % 12
	if hour12 == 0 {
		hour12 = 12
	}

	// Ordinal suffix (st, nd, rd, th)
	ordinal := func(n int) string {
		v := n % 100
		s := "th"
		if v < 20 {
			switch v {
			case 1:
				s = "st"
			case 2:
				s = "nd"
			case 3:
				s = "rd"
			}
		} else {
			switch v % 10 {
			case 1:
				s = "st"
			case 2:
				s = "nd"
			case 3:
				s = "rd"
			}
		}
		return fmt.Sprintf("%d%s", n, s)
	}

	return fmt.Sprintf("%s %s of %s %d %02d:%02d:%02d %s",
		dayName, ordinal(day), month, year, hour12, minute, second, period)
}
