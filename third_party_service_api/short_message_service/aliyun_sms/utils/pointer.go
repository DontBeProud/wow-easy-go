package utils

import (
	"strconv"
	"time"
)

// ParseStrPointerIntoString *string -> string. *string等于nil时返回""
func ParseStrPointerIntoString(raw *string) string {
	if raw != nil {
		return *raw
	}
	return ""
}

// ParseStrPointerIntoInt64 *string -> int64. *string等于nil时返回0
func ParseStrPointerIntoInt64(raw *string) (int, error) {
	if raw == nil {
		return 0, nil
	}
	return strconv.Atoi(*raw)
}

// ParseIntPointerIntoInt *int64 -> int64. *int64等于nil时返回0
func ParseIntPointerIntoInt(raw *int64) int64 {
	if raw == nil {
		return 0
	}
	return *raw
}

// ParseStrPointerIntoTime *string -> time.Time. *string等于nil时返回空时间
func ParseStrPointerIntoTime(raw *string, dateFmt string) time.Time {
	if raw == nil {
		return time.Time{}
	}

	t, err := time.ParseInLocation(dateFmt, *raw, time.Local)
	if err != nil {
		return time.Time{}
	}

	return t
}
