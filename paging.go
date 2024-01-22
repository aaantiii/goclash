package clash

// Paging represents the paging information returned by the API.
type Paging struct {
	Cursors PagingCursors `json:"cursors,omitempty"`
}

// PagingCursors represents the paging cursors returned by the API.
type PagingCursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

// PagingParams represents the parameters for a paginated request.
type PagingParams struct {
	PagingCursors
	Limit int `json:"limit,omitempty"`
}

// PaginatedResponse represents a paginated response from the API.
type PaginatedResponse[T any] struct {
	Paging Paging `json:"paging,omitempty"`
	Items  []T    `json:"items,omitempty"`
}
