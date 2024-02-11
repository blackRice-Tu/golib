package xcommon

func ToPtr[T any](item T) *T {
	return &item
}

func ToPtrs[T any](items []T) []*T {
	ptrs := make([]*T, len(items))
	for i := range items {
		ptrs[i] = &items[i]
	}
	return ptrs
}
