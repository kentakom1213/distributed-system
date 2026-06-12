package storage

import "time"

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func nowJST() time.Time {
	return time.Now().In(jst)
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, s)
}
