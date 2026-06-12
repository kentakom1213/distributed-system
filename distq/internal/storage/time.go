package storage

import "time"

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func nowJST() time.Time {
	return time.Now().In(jst)
}
