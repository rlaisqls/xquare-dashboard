package query

import (
	"github.com/timberio/go-datemath"
	"strconv"
	"time"
)

type DataTimeRange struct {
	From string
	To   string
	Now  time.Time
}

func (tr DataTimeRange) GetFromAsTimeUTC() time.Time {
	return tr.MustGetFrom().UTC()
}

func (tr DataTimeRange) GetToAsTimeUTC() time.Time {
	return tr.MustGetTo().UTC()
}

func (tr DataTimeRange) MustGetFrom() time.Time {
	res, err := tr.ParseFrom()
	if err != nil {
		return time.Unix(0, 0)
	}
	return res
}

func (tr DataTimeRange) MustGetTo() time.Time {
	res, err := tr.ParseTo()
	if err != nil {
		return time.Unix(0, 0)
	}
	return res
}

func (tr DataTimeRange) ParseFrom() (time.Time, error) {
	pt := newParsableTime(tr.From)
	return pt.Parse()
}

func (tr DataTimeRange) ParseTo() (time.Time, error) {
	pt := newParsableTime(tr.To)
	return pt.Parse()
}

func (t parsableTime) Parse() (time.Time, error) {
	// Milliseconds since Unix epoch.
	if val, err := strconv.ParseInt(t.time, 10, 64); err == nil {
		return time.UnixMilli(val), nil
	}

	// Duration relative to current time.
	if diff, err := time.ParseDuration("-" + t.time); err == nil {
		return t.now.Add(diff), nil
	}

	// Advanced time string, mimics the frontend's datemath library.
	return datemath.ParseAndEvaluate(t.time, t.datemathOptions()...)
}

func (t parsableTime) datemathOptions() []func(*datemath.Options) {
	options := []func(*datemath.Options){
		datemath.WithNow(t.now),
	}
	if t.location != nil {
		options = append(options, datemath.WithLocation(t.location))
	}
	if t.weekstart != nil {
		options = append(options, datemath.WithStartOfWeek(*t.weekstart))
	}
	return options
}

type parsableTime struct {
	time             string
	now              time.Time
	location         *time.Location
	weekstart        *time.Weekday
	fiscalStartMonth *time.Month
	roundUp          bool
}

func newParsableTime(t string) parsableTime {
	return parsableTime{
		time: t,
		now:  time.Now(),
	}
}
