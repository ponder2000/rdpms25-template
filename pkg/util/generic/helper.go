package generic

import (
	"encoding/json"
	"errors"
	"reflect"
)

func InLine[T any](val bool, v1 T, v2 T) T {
	if val {
		return v1
	}
	return v2
}

// DeepCopyWithJSON : T should be a pointer
func DeepCopyWithJSON[T any](src T, dest T) error {
	raw, e := json.Marshal(src)
	if e != nil {
		return e
	}

	return json.Unmarshal(raw, dest)
}

func RemoveDuplicate[T comparable](arr []T) []T {
	v := map[T]struct{}{}
	res := make([]T, 0)
	for _, a := range arr {
		_, ok := v[a]
		if !ok {
			res = append(res, a)
		}
		v[a] = struct{}{}
	}
	return res
}

func Mapper[T any, V any](arr []T, extractor func(T) V) []V {
	res := make([]V, len(arr))
	for i, a := range arr {
		res[i] = extractor(a)
	}
	return res
}

func Transform[T1 any, T2 any](arr []T1, convert func(T1) T2) []T2 {
	res := make([]T2, 0)
	for _, a := range arr {
		res = append(res, convert(a))
	}
	return res
}

func Where[T any](arr []T, isAcceptable func(T) bool) []T {
	res := make([]T, 0)
	for i, a := range arr {
		if isAcceptable(arr[i]) {
			res = append(res, a)
		}
	}
	return res
}

func WhereFirst[T any](arr []T, isAcceptable func(T) bool) (T, error) {
	for i, a := range arr {
		if isAcceptable(arr[i]) {
			return a, nil
		}
	}
	var zero T
    return zero, errors.New("no matching element found")
}

func WhereIndex[T any](arr []T, isAcceptableIndex func(int) bool) []T {
	res := make([]T, 0)
	for i, a := range arr {
		if isAcceptableIndex(i) {
			res = append(res, a)
		}
	}
	return res
}

func ConvertSliceToMap[K comparable, V any](arr []V, key func(V) K) map[K]V {
	res := make(map[K]V)
	for i := range arr {
		res[key(arr[i])] = arr[i]
	}
	return res
}

func GetValuesFromMap[T comparable, V any](m map[T]V) []V {
	res := make([]V, 0)
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func GetKeysFromMap[T comparable, V any](m map[T]V) []T {
	res := make([]T, 0)
	for k := range m {
		res = append(res, k)
	}
	return res
}

func UnionSlice[T any, K comparable](arr1 []T, arr2 []T, comparator func(T) K) []T {
	visited := make(map[K]T)
	for _, a := range arr1 {
		visited[comparator(a)] = a
	}
	for _, a := range arr2 {
		visited[comparator(a)] = a
	}
	return GetValuesFromMap(visited)
}

func IntersectionSlice[T any, K comparable](arr1 []T, arr2 []T, comparator func(T) K) []T {
	res := make([]T, 0)

	visited := make(map[K]T)
	for _, a := range arr1 {
		visited[comparator(a)] = a
	}
	for _, a := range arr2 {
		val, ok := visited[comparator(a)]
		if ok {
			res = append(res, val)
		}
	}
	return res
}

func FlattenMap(nestedMap map[string]any) map[string]any {
	flatMap := make(map[string]any)
	flattenHelper(nestedMap, flatMap)
	return flatMap
}

func flattenHelper(nestedMap any, flatMap map[string]any) {
	switch reflect.TypeOf(nestedMap).Kind() {
	case reflect.Map:
		for k, v := range nestedMap.(map[string]any) {
			if v != nil && reflect.TypeOf(v).Kind() == reflect.Map {
				flattenHelper(v, flatMap)
			} else {
				flatMap[k] = v
			}
		}
	default:
		// Handle non-map values (if necessary)
	}
}
