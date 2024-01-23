package goclash

import (
	"net/http"

	"github.com/bytedance/sonic"
)

type LabelsData struct {
	Paging *Paging `json:"paging,omitempty"`
	Labels []Label `json:"items,omitempty"`
}

type Label struct {
	IconUrls ImageURLs `json:"iconUrls,omitempty"`
	Name     string    `json:"name,omitempty"`
	ID       int       `json:"id,omitempty"`
}

// GetPlayerLabels returns a paginated list of player labels. Pass params=nil to get all labels.
func (h *Client) GetPlayerLabels(params *PagingParams) (*PaginatedResponse[Label], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LabelsEndpoint.Build("players"), req, true)
	if err != nil {
		return nil, err
	}

	var labels *PaginatedResponse[Label]
	err = sonic.Unmarshal(data, &labels)
	return labels, err
}

// GetClanLabels returns a paginated list of clan labels. Pass params=nil to get all labels.
func (h *Client) GetClanLabels(params *PagingParams) (*PaginatedResponse[Label], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LabelsEndpoint.Build("clans"), req, true)
	if err != nil {
		return nil, err
	}

	var labels *PaginatedResponse[Label]
	err = sonic.Unmarshal(data, &labels)
	return labels, err
}
