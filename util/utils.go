package util

import "time"

var inputFormatExample = "15:04:05-07:00"
var outputFormatExample = "15:04:05"

func NormalizeTimeToUTC(tStr string) (string, error) {
	parsed, err := time.Parse(inputFormatExample, tStr)
	if err != nil {
		// fallback to plain time, assume UTC
		parsed, err = time.Parse("15:04:05", tStr)
		if err != nil {
			return "", err
		}
	}

	utc := parsed.UTC()

	return utc.Format(outputFormatExample), nil
}
