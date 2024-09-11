package util

import (
	"time"
)

func ParseDateTime(dateString string) (time.Time, error) {
	const layout = "2006-01-02 15:04:05"
	t, err := time.Parse(layout, dateString)
	return t, err
}

func WibToUtc(wib time.Time) time.Time {
	return wib.Add(-7 * time.Hour)
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func StdTimeToString(t *time.Time, convert bool) string {
	result := time.Time{}
	if t != nil {
		result = *t
	}

	if convert {
		result = WibToUtc(result)
	}

	return result.Format("2006-01-02T15:04:05+00:00")
}

func DateToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func TimeToStandardUtc(dateString string) (string, error) {
	parsedTime, err := ParseDateTime(dateString)
	if err != nil {
		return "", err
	}
	standardUtc := TimeToString(WibToUtc(parsedTime))
	return standardUtc, nil
}

func TimeToStandard(dateString string) (string, error) {
	parsedTime, err := ParseDateTime(dateString)
	if err != nil {
		return "", err
	}
	standardUtc := TimeToString(parsedTime)
	return standardUtc, nil
}

func TimeConvert(dateString string, convertToUtc bool) string {
	var result string
	var err error

	if convertToUtc {
		result, err = TimeToStandardUtc(dateString)
	} else {
		result, err = TimeToStandard(dateString)
	}

	if err != nil {
		return ""
	}

	return result

}

func DateConvert(dateString string, convertToUtc bool) string {

	parsedTime, err := ParseDateTime(dateString)
	if err != nil {
		return ""
	}

	var result string
	if convertToUtc {
		result = DateToString(WibToUtc(parsedTime))
	} else {
		result = DateToString(parsedTime)
	}

	return result

}
