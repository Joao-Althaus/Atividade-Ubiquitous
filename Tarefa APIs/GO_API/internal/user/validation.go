package user

import "time"

func parseBirthday(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}
