package data

import "time"

func CorrectNotitifactionTime(check time.Time, cake *Cake) bool {

	t, _ := time.Parse(time.RFC3339, cake.When)

	if time.Until(t) > (time.Hour * 24 * 5) {
		return false
	}

	init := time.Now()
	start := time.Date(init.Year(), init.Month(), init.Day(), 9, 0, 0, 0, time.UTC)
	end := time.Date(init.Year(), init.Month(), init.Day(), 23, 0, 0, 0, time.UTC)

	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}
