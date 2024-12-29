package models

type PaginationResponse[T Source | Reservation | Customer] struct {
	Total   int64
	Page    uint32
	Count   int
	Content []T
}

func NewPaginationResponse[T Source | Reservation | Customer](vals []T, total int64, page uint32) PaginationResponse[T] {
	return PaginationResponse[T]{
		Content: vals,
		Page:    page,
		Total:   total,
		Count:   len(vals),
	}
}
