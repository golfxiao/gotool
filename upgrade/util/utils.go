package util

import (
	"os"
	"strconv"

	"github.com/sluu99/uuid"
)

func CreateDir(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return false
			}
			return true
		}
	}
	return true
}

func ToString(v interface{}) string {
	if r, ok := v.(string); ok {
		return r
	}
	//not string should be convert number to string
	switch v.(type) {
	case uint64:
		return strconv.Itoa(int(v.(uint64)))
	case int64:
		return strconv.Itoa(int(v.(int64)))
	case int:
		return strconv.Itoa((v.(int)))
	case int32:
		return strconv.Itoa(int(v.(int32)))
	case uint32:
		return strconv.Itoa(int(v.(uint32)))
	case float64:
		return strconv.Itoa(int(v.(float64)))
	case int8:
		return strconv.Itoa(int(v.(int8)))
	case uint8:
		return strconv.Itoa(int(v.(uint8)))
	case bool:
		if v.(bool) {
			return "true"
		} else {
			return "false"
		}
	}
	return ""
}

func ToInt64(v interface{}, defaultVal int64) int64 {
	if v == nil {
		return defaultVal
	}

	switch v.(type) {
	case bool:
		if v.(bool) {
			return 1
		}
		return 0
	case string:
		i, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return defaultVal
		}
		return i
	case uint64:
		return int64(v.(uint64))
	case int64:
		return int64(v.(int64))
	case int:
		return int64(v.(int))
	case int32:
		return int64(v.(int32))
	case uint32:
		return int64(v.(uint32))
	case float64:
		return int64(v.(float64))
	case int8:
		return int64(v.(int8))
	case uint8:
		return int64(v.(uint8))
	}
	return defaultVal
}

func UUID() string {
	// uuid.Rand().Hex()
	return uuid.Rand().Hex()
}
