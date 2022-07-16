package formats

import "time"

type DateTime struct {
	time.Time
}

type Amount struct {
	amount float64
}
