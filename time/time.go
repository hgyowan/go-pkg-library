package time

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

const (
	CRITERIAL_TIME_LAYOUT = "2006-01-02"
	DEFAULT_TIME_LAYOUT   = "2006-01-02 15:04:05"
	PARSE_TIME_LAYOUT     = "2006-01-02T15:04:05"
)

const (
	TIME_UNIT_DAY   = "DAY"
	TIME_UNIT_MONTH = "MONTH"
	TIME_UNIT_YEAR  = "YEAR"
)

func ConvertToSafeTimestamp[T any](input T) *timestamppb.Timestamp {
	switch v := any(input).(type) {
	case *time.Time:
		if v == nil {
			return nil
		}
		return timestamppb.New(*v)
	case time.Time:
		if v.IsZero() {
			return nil
		}
		return timestamppb.New(v)
	case gorm.DeletedAt:
		if v.Time.IsZero() {
			return nil
		}
		return timestamppb.New(v.Time)
	default:
		return nil
	}
}

func CustomAddDate(date time.Time, years int, months int, days int) time.Time {
	if months != 0 {
		// date 의 날짜가 마지막 날이면 months 를 더했을때 다음달도 말일로 처리되어야함.
		if isLastDayOfMonth(date) {
			year := date.Year()
			month := date.Month() + time.Month(months)
			day := getLastDayOfNextMonth(date, months)
			hour := date.Hour()
			minute := date.Minute()
			second := date.Second()

			date = time.Date(year, month, day, hour, minute, second, 0, time.UTC)
		} else {
			date = date.AddDate(0, months, 0)
		}
	}

	if years != 0 {
		date = date.AddDate(years, 0, 0)
	}

	if days != 0 {
		date = date.AddDate(0, 0, days)
	}

	return date
}

func isLastDayOfMonth(date time.Time) bool {
	// 현재 날짜에서 하루를 더한 날짜
	nextDay := date.AddDate(0, 0, 1)
	// 다음 날의 월과 현재 날의 월이 다르면 현재 날은 그 달의 마지막 날
	return date.Month() != nextDay.Month()
}

func getLastDayOfNextMonth(date time.Time, addedMonth int) int {
	year := date.Year()
	month := date.Month() + time.Month(addedMonth) + 1
	// 다다음달 1일 에서 -1 일
	nextMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := nextMonth.AddDate(0, 0, -1) // 하루를 빼면 현재 달의 마지막 날
	return lastDay.Day()
}

func ConvertTimeStampPbToTime(pbTime *timestamppb.Timestamp) *time.Time {
	var t *time.Time
	if pbTime != nil {
		tempTime := pbTime.AsTime()
		t = &tempTime
	}

	return t
}
