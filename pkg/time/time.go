package time

import "time"

type TuiterTime struct {
	location *time.Location
}

func (t *TuiterTime) Now() time.Time {
	return time.Now().In(t.location)
}

func NewTuiterTime(location *time.Location) *TuiterTime {
	return &TuiterTime{location: location}
}

type Time interface {
	Now() time.Time
}
