package statistical

import "time"

const (
	dateLayout      = "2006-01-02"
	dateShortLayout = "20060102"
)

func GetTimeStr(t time.Time) string {
	return t.Format(dateLayout)
}

func GetShortTimeStr(t time.Time) string {
	return t.Format(dateShortLayout)
}
