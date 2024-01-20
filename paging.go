package clash

type Paging struct {
	Cursors PagingCursors `json:"cursors,omitempty"`
}

type PagingCursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}
