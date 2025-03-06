package utils

func NewListResult[T any](data T, total int64) *ListPublicResponse[T] {
	return &ListPublicResponse[T]{Data: data, Total: total}
}

type ListPublicResponse[T any] struct {
	Data  T     `json:"data"`
	Total int64 `json:"total,omitempty"`
}
