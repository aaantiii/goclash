package clash

type LabelsData struct {
	Paging *Paging `json:"paging,omitempty"`
	Labels []Label `json:"items,omitempty"`
}

type Label struct {
	IconUrls ImageURLs `json:"iconUrls,omitempty"`
	Name     string    `json:"name,omitempty"`
	ID       int       `json:"id,omitempty"`
}
