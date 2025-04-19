package parser

import "encoding/json"

func MapToJson[K comparable, T any](v map[K]T) ([]byte, error) {
	return json.Marshal(v)
}

func MapToJsonWithoutError[K comparable, T any](v map[K]T) []byte {
	res, _ := MapToJson(v)
	return res
}

func JsonToMap[K comparable, T any](v []byte) (map[K]T, error) {
	res := make(map[K]T)
	e := json.Unmarshal(v, &res)
	return res, e
}

func JsonToMapWithoutError[K comparable, T any](v []byte) map[K]T {
	res, _ := JsonToMap[K, T](v)
	return res
}

func ArrayToJson[T any](v []T) ([]byte, error) {
	return json.Marshal(v)
}

func ArrayToJsonWithoutError[T any](v []T) []byte {
	res, _ := ArrayToJson[T](v)
	return res
}

func JsonToArray[T any](v []byte) ([]T, error) {
	res := make([]T, 0)
	e := json.Unmarshal(v, &res)
	return res, e
}

func JsonToArrayWithoutError[T any](v []byte) []T {
	res, _ := JsonToArray[T](v)
	return res
}
