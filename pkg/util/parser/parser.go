package parser

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

func ExtractRequiredString(rawVal string) (string, error) {
	if len(rawVal) > 0 {
		return rawVal, nil
	}
	return "", errors.New("required field error")
}

func ExtractInt(rawVal string, required bool) (int, error) {
	var val int
	if len(rawVal) > 0 {
		var _e error
		if val, _e = strconv.Atoi(rawVal); _e != nil {
			return 0, errors.New("invalid int")
		}
	} else if required {
		return 0, errors.New("required field error")
	}
	return val, nil
}

func ExtractInt64(rawVal string, required bool) (int64, error) {
	var val int64
	if len(rawVal) > 0 {
		var _e error
		if val, _e = strconv.ParseInt(rawVal, 10, 64); _e != nil {
			return 0, errors.New("invalid int64")
		}
	} else if required {
		return 0, errors.New("required int64 error")
	}
	return val, nil
}

func ExtractFloat(rawVal string, required bool) (float64, error) {
	var val float64
	if len(rawVal) > 0 {
		var _e error
		if val, _e = strconv.ParseFloat(rawVal, 64); _e != nil {
			return 0, errors.New("invalid float")
		}
	} else if required {
		return 0, errors.New("required float64 error")
	}
	return val, nil
}

func ExtractBool(rawVal string) (bool, error) {
	if len(rawVal) > 0 {
		if rawVal == "true" {
			return true, nil
		}
		if rawVal == "false" {
			return false, nil
		}
	}
	return false, errors.New("invalid boolean")
}

func ExtractIntSlice(rawVal string, required bool) ([]int, error) {
	if required && len(rawVal) == 0 {
		return nil, errors.New("required int array error")
	} else if !required && len(rawVal) == 0 {
		return nil, nil
	}

	res := make([]int, 0)
	stringSlice := strings.Split(rawVal, ",")

	for _, s := range stringSlice {
		i, e := strconv.Atoi(s)
		if e != nil {
			return res, errors.New("invalid int array passed!")
		} else {
			res = append(res, i)
		}
	}

	return res, nil
}

func ExtractTimeFromStringEpoch(rawVal string) (time.Time, error) {
	epoch, e := ExtractInt(rawVal, true)
	if e != nil {
		return time.Time{}, e
	}

	// epoch is divided with 1000 assuming the ints will be pased in milliseconds
	return time.Unix(int64(epoch/1000), 0), e
}

func TypeToMap[T any](obj T) (map[string]any, error) {
	raw, e := json.Marshal(obj)
	if e != nil {
		return nil, e
	}

	res := make(map[string]any)
	if e := json.Unmarshal(raw, &res); e != nil {
		return nil, e
	}
	return res, nil
}

func TypeToMapWithBlacklistKeys[T any](obj T, deleteKey ...string) (map[string]any, error) {
	mapObj, e := TypeToMap(obj)
	if e != nil {
		return nil, e
	}
	for _, k := range deleteKey {
		delete(mapObj, k)
	}
	return mapObj, nil
}
