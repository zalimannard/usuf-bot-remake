package util

import (
	"fmt"
	"time"
)

func FormatAsHHMMSS(d time.Duration) string {
	d = d.Round(time.Second)

	totalSec := int64(d / time.Second)
	if totalSec < 0 {
		totalSec = -totalSec
	}
	h := totalSec / 3600
	m := (totalSec % 3600) / 60
	s := totalSec % 60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
