package goclash

import (
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
)

type Location struct {
	LocalizedName string `json:"localizedName,omitempty"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	IsCountry     bool   `json:"isCountry"`
	CountryCode   string `json:"countryCode"`
}

type ClanRanking struct {
	ClanLevel    int       `json:"clanLevel"`
	ClanPoints   int       `json:"clanPoints"`
	Location     Location  `json:"location"`
	Members      int       `json:"members"`
	Tag          string    `json:"tag"`
	Name         string    `json:"name"`
	Rank         int       `json:"rank"`
	PreviousRank int       `json:"previousRank"`
	BadgeURLs    ImageURLs `json:"badgeUrls"`
}

type PlayerRanking struct {
	League       League            `json:"league"`
	Clan         PlayerRankingClan `json:"clan"`
	AttackWins   int               `json:"attackWins"`
	DefenseWins  int               `json:"defenseWins"`
	Tag          string            `json:"tag"`
	Name         string            `json:"name"`
	ExpLevel     int               `json:"expLevel"`
	Rank         int               `json:"rank"`
	PreviousRank int               `json:"previousRank"`
	Trophies     int               `json:"trophies"`
}

type PlayerBuilderBaseRanking struct {
	BuilderBaseLeague   BuilderBaseLeague `json:"builderBaseLeague"`
	Clan                PlayerRankingClan `json:"clan"`
	Tag                 string            `json:"tag"`
	Name                string            `json:"name"`
	ExpLevel            int               `json:"expLevel"`
	Rank                int               `json:"rank"`
	PreviousRank        int               `json:"previousRank"`
	BuilderBaseTrophies int               `json:"builderBaseTrophies"`
}

type ClanBuilderBaseRanking struct {
	ClanPoints            int `json:"clanPoints"`
	ClanBuilderBasePoints int `json:"clanBuilderBasePoints"`
}

type ClanCapitalRanking struct {
	ClanPoints        int `json:"clanPoints"`
	ClanCapitalPoints int `json:"clanCapitalPoints"`
}

// GetClanRankings returns a paginated list of clan rankings for a specific location.
//
// GET /locations/{locationId}/rankings/clans
func (h *Client) GetClanRankings(locationID int, params *PagingParams) (*PaginatedResponse[ClanRanking], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LocationsEndpoint.Build(strconv.Itoa(locationID), "rankings/clans"), req, true)
	if err != nil {
		return nil, err
	}

	var rankings *PaginatedResponse[ClanRanking]
	err = sonic.Unmarshal(data, &rankings)
	return rankings, err
}

// GetPlayerRankings returns a paginated list of player rankings for a specific location.
//
// GET /locations/{locationId}/rankings/players
func (h *Client) GetPlayerRankings(locationID int, params *PagingParams) (*PaginatedResponse[PlayerRanking], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LocationsEndpoint.Build(strconv.Itoa(locationID), "rankings/players"), req, true)
	if err != nil {
		return nil, err
	}

	var rankings *PaginatedResponse[PlayerRanking]
	err = sonic.Unmarshal(data, &rankings)
	return rankings, err
}

// GetPlayerBuilderBaseRankings returns a paginated list of player builder base rankings for a specific location.
//
// GET /locations/{locationId}/rankings/players-builder-base
func (h *Client) GetPlayerBuilderBaseRankings(locationID int, params *PagingParams) (*PaginatedResponse[PlayerBuilderBaseRanking], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LocationsEndpoint.Build(strconv.Itoa(locationID), "rankings/players-builder-base"), req, true)
	if err != nil {
		return nil, err
	}

	var rankings *PaginatedResponse[PlayerBuilderBaseRanking]
	err = sonic.Unmarshal(data, &rankings)
	return rankings, err
}

// GetClanBuilderBaseRankings returns a paginated list of clan builder base rankings for a specific location.
//
// GET /locations/{locationId}/rankings/clans-builder-base
func (h *Client) GetClanBuilderBaseRankings(locationID int, params *PagingParams) (*PaginatedResponse[ClanBuilderBaseRanking], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LocationsEndpoint.Build(strconv.Itoa(locationID), "rankings/clans-builder-base"), req, true)
	if err != nil {
		return nil, err
	}

	var rankings *PaginatedResponse[ClanBuilderBaseRanking]
	err = sonic.Unmarshal(data, &rankings)
	return rankings, err
}

// GetLocations returns a paginated list of all available locations.
//
// GET /locations
func (h *Client) GetLocations(params *PagingParams) (*PaginatedResponse[Location], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LocationsEndpoint.Build(), req, true)
	if err != nil {
		return nil, err
	}

	var locations *PaginatedResponse[Location]
	err = sonic.Unmarshal(data, &locations)
	return locations, err
}

// GetClanCapitalRankings returns a paginated list of clan capital rankings for a specific location.
//
// GET /locations/{locationId}/rankings/capitals
func (h *Client) GetClanCapitalRankings(locationID int, params *PagingParams) (*PaginatedResponse[ClanCapitalRanking], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LocationsEndpoint.Build(strconv.Itoa(locationID), "rankings/capitals"), req, true)
	if err != nil {
		return nil, err
	}

	var rankings *PaginatedResponse[ClanCapitalRanking]
	err = sonic.Unmarshal(data, &rankings)
	return rankings, err
}

// GetLocation returns information about a specific location.
//
// GET /locations/{locationId}
func (h *Client) GetLocation(locationID int) (*Location, error) {
	data, err := h.do(http.MethodGet, LocationsEndpoint.Build(strconv.Itoa(locationID)), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var location *Location
	err = sonic.Unmarshal(data, &location)
	return location, err
}
