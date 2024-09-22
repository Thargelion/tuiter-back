package instant

import "time"

type LocatedInstant struct {
	location *time.Location
}

func (t *LocatedInstant) Now() time.Time {
	return time.Now().In(t.location)
}

func NewTuiterTime(location *time.Location) *LocatedInstant {
	return &LocatedInstant{location: location}
}

type Instant interface {
	Now() time.Time
}
