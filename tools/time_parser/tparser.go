package time_parser

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	DATE = "2006-01-02"
)

func GetDuration(val string) string {
	start := strings.Index(val, "T")
	end := strings.Index(val, "Z")
	return val[start+1 : end]
}

func StringToTime(val string) time.Time {
	date, _ := time.Parse(DATE, val)
	return date
}

func TrackDuration(path string) (string, error) {
	out, err := exec.Command("sox", "-t", "mp3", path, "-n", "stat").CombinedOutput()
	if err != nil {
		return "", err
	}

	info := strings.Split(string(out), "\n")[1]
	k := strings.TrimSpace(strings.Split(info, ":")[1])

	t, err := strconv.ParseFloat(k, 64)
	if err != nil {
		return "", err
	}

	hours := uint(t) / 3600
	minutes := (uint(t) - (3600 * hours)) / 60
	seconds := uint(t) - (3600 * hours) - (minutes * 60)

	return fmt.Sprintf("%d:%d:%d", hours, minutes, seconds), nil
}
