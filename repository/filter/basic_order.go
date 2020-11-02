package filter

// swagger:parameters basicOrder
type BasicOrder struct {
	// Field which you want to order by
	OrderBy []string `json:"order_by" schema:"order_by"`
}

func NewBasicOrder() *BasicOrder {
	return &BasicOrder{}
}

func (s *BasicOrder) GetOrderBy() []string {
	return s.OrderBy
}
