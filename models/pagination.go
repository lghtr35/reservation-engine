package models

type Pagination struct {
	Page uint32 `json:"page"`
	Size uint32 `json:"size"`
}

func (p *Pagination) Offset() int {
	return (int(p.Page) - 1) * int(p.Size)
}
