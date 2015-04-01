package codebase

import (
    "math"
    "time"
)

type WorkingDays time.Time

func (w WorkingDays) time() time.Time {
    return time.Time(w)
}

func (w WorkingDays) SinceUnix() int64 {
    beginningOfDay := w.time().UTC().Truncate(24 * time.Hour)

    days := beginningOfDay.UnixNano() / int64(24*time.Hour)
    days += 3 // Move day to Sunday

    working_days := float64((days / 7) * 5)
    working_days += math.Min(float64(days%7), 5)

    return int64(working_days - 3) // Move back
}
