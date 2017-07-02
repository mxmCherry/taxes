package taxes

import "time"

// Quarter holds financial year quarter.
type Quarter struct {
	// Year holds quarter's year.
	Year int
	// Quarter holds quarter's number (1-based).
	Quarter int
}

// QuarterOf returns quarted of a given time.
func QuarterOf(t time.Time) Quarter {
	y, m, _ := t.Date()
	return Quarter{
		Year:    y,
		Quarter: int((m-1)/3 + 1),
	}
}
