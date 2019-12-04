package time_parser

import (
	"strings"
)

func GetDuration(val string) string {
	start := strings.Index(val, "T")
	end := strings.Index(val, "Z")
	return val[start+1 : end]
}
