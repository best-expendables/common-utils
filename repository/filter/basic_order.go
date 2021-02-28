package filter

type BasicOrder struct {
	OrderBy []string `json:"orderBy"`
}

func NewBasicOrder() *BasicOrder {
	return &BasicOrder{}
}

func (s *BasicOrder) GetOrderBy() []string {
	return s.OrderBy
}
