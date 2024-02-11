package xalgorithm

import (
	"golang.org/x/exp/rand"
	"time"
)

// SubSlice like Python array
func SubSlice[T any](slice []T, args ...int) []T {
	length := len(slice)
	if length == 0 {
		return []T{}
	}
	if len(args) == 0 {
		return slice
	}

	begin := 0
	end := 0
	if len(args) == 1 { // only with begin
		begin = args[0]
		end = length
	}

	if len(args) == 2 {
		begin = args[0]
		end = args[1]
	}

	if begin >= length {
		return []T{}
	}

	if begin < 0 {
		begin = length + begin
	}
	if end < 0 {
		end = length + end
	}
	if begin >= end {
		return []T{}
	}
	return slice[begin:end]
}

func SliceChunk[T any](items []T, size int) (chunks [][]T) {
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}
	return append(chunks, items)
}

// IsUnited to check if a team is united
func IsUnited[T comparable](goal T, options ...T) bool {
	if len(options) == 0 {
		return true
	}
	selected := options[0] == goal
	for _, value := range options[1:] {
		if selected != (value == goal) {
			return false
		}
	}
	return true
}

func SliceRandom[T any](items []T, size int) []T {
	length := len(items)
	if size < 0 {
		size = 0
	} else if size > length {
		size = length
	}
	tmpItems := make([]T, length)
	copy(tmpItems, items)
	rand.Seed(uint64(time.Now().UnixNano()))
	rand.Shuffle(len(tmpItems), func(i, j int) {
		tmpItems[i], tmpItems[j] = tmpItems[j], tmpItems[i]
	})
	return tmpItems[:size]
}
