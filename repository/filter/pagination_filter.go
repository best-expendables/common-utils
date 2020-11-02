package filter

const (
	defaultPerPage = 50
	maxPerPage     = 200
)

type PaginationFilter struct {
	BasicFilter
	BasicOrder `swaggerignore:"true"`

	Page    int `json:"page" schema:"page"`
	PerPage int `json:"per_page" schema:"per_page"`
}

func NewPaginationFilter() *PaginationFilter {
	return &PaginationFilter{
		BasicFilter: *NewBasicFilter(),
		BasicOrder:  *NewBasicOrder(),
	}
}

// implement repository.Filter interface
func (f *PaginationFilter) GetLimit() int {
	return f.GetPerPage() + 1
}

// implement repository.Filter interface
func (f *PaginationFilter) GetOffset() int {
	return (f.GetPage() - 1) * f.GetPerPage()
}

func (f *PaginationFilter) GetPage() int {
	if f.Page < 1 {
		return 1
	}
	return f.Page
}

func (f *PaginationFilter) GetPerPage() int {
	if f.PerPage < 1 || f.PerPage > maxPerPage {
		return defaultPerPage
	}

	return f.PerPage
}
