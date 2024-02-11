package xcommon

import (
	"fmt"
	"strings"
)

type Numeric interface {
	int | int32 | int64 | float64 | float32
}

func SliceToSet[T comparable](items []T) map[T]struct{} {
	mapping := make(map[T]struct{})
	for _, item := range items {
		mapping[item] = struct{}{}
	}
	return mapping
}

func IntSlice[T Numeric](items []T) []int {
	slice := make([]int, len(items))
	for i := range items {
		slice[i] = int(items[i])
	}
	return slice
}

func JoinSlice[T any](items []T, sep string) string {
	stringItems := make([]string, len(items))
	for i, item := range items {
		stringItems[i] = fmt.Sprintf("%v", item)
	}
	return strings.Join(stringItems, sep)
}
