package clash

import (
	"net/http"

	"github.com/bytedance/sonic"
)

type GoldPassSeason struct {
	StartTime string
	EndTime   string
}

// GetCurrentGoldPassSeason returns the current gold pass season.
//
// GET /goldpass/seasons/current
func (h *Client) GetCurrentGoldPassSeason() (*GoldPassSeason, error) {
	req := h.withAuth(h.newDefaultRequest())
	data, err := h.do(http.MethodGet, GoldPassEndpoint.Build(), req, true)
	if err != nil {
		return nil, err
	}

	var season *GoldPassSeason
	err = sonic.Unmarshal(data, &season)
	return season, err
}
