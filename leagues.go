package goclash

import (
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
)

type WarLeague struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// CWL leagues
const (
	WarLeagueUnranked = 48000000 + iota
	WarLeagueBronzeIII
	WarLeagueBronzeII
	WarLeagueBronzeI
	WarLeagueSilverIII
	WarLeagueSilverII
	WarLeagueSilverI
	WarLeagueGoldIII
	WarLeagueGoldII
	WarLeagueGoldI
	WarLeagueCrystalIII
	WarLeagueCrystalII
	WarLeagueCrystalI
	WarLeagueMasterIII
	WarLeagueMasterII
	WarLeagueMasterI
	WarLeagueChampionIII
	WarLeagueChampionII
	WarLeagueChampionI
)

// Home village leagues
const (
	LeagueUnranked = 29000000 + iota
	LeagueBronzeIII
	LeagueBronzeII
	LeagueBronzeI
	LeagueSilverIII
	LeagueSilverII
	LeagueSilverI
	LeagueGoldIII
	LeagueGoldII
	LeagueGoldI
	LeagueCrystalIII
	LeagueCrystalII
	LeagueCrystalI
	LeagueMasterIII
	LeagueMasterII
	LeagueMasterI
	LeagueChampionIII
	LeagueChampionII
	LeagueChampionI
	LeagueTitanIII
	LeagueTitanII
	LeagueTitanI
	LeagueLegend
)

type LeagueData struct {
	Paging  Paging   `json:"paging,omitempty"`
	Leagues []League `json:"items,omitempty"`
}

type CapitalLeague struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type League struct {
	IconUrls ImageURLs `json:"iconUrls,omitempty"`
	Name     string    `json:"name,omitempty"`
	ID       int       `json:"id,omitempty"`
}

// PlayerRankingList contains information about a player's ranking.
type PlayerRankingList struct {
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

type PlayerRankingClan struct {
	Tag       string    `json:"tag"`
	Name      string    `json:"name"`
	BadgeURLs ImageURLs `json:"badgeUrls"`
}

type LeagueSeason struct {
	ID string `json:"id"`
}

// GetCapitalLeagues returns a paginated list of capital leagues. Pass params=nil to get all leagues.
//
// GET /capitalleagues
func (h *Client) GetCapitalLeagues(params *PagingParams) (*PaginatedResponse[CapitalLeague], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, CapitalLeaguesEndpoint.Build(), req, true)
	if err != nil {
		return nil, err
	}

	var leagues *PaginatedResponse[CapitalLeague]
	err = sonic.Unmarshal(data, &leagues)
	return leagues, err
}

// GetLeagues returns a paginated list of leagues. Pass params=nil to get all leagues.
//
// GET /leagues
func (h *Client) GetLeagues(params *PagingParams) (*PaginatedResponse[League], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LeaguesEndpoint.Build(), req, true)
	if err != nil {
		return nil, err
	}

	var leagues *PaginatedResponse[League]
	err = sonic.Unmarshal(data, &leagues)
	return leagues, err
}

// GetLegendLeagueRanking returns a paginated list of players in the provided legend league season.
//
// GET /leagues/{leagueId}/seasons/{seasonId}
func (h *Client) GetLegendLeagueRanking(leagueID, seasonID string, params *PagingParams) (*PaginatedResponse[PlayerRankingList], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LeaguesEndpoint.Build(leagueID, "seasons", seasonID), req, true)
	if err != nil {
		return nil, err
	}

	var rankings *PaginatedResponse[PlayerRankingList]
	err = sonic.Unmarshal(data, &rankings)
	return rankings, err
}

// GetCapitalLeague returns information about a single capital league.
//
// GET /capitalleagues/{leagueId}
func (h *Client) GetCapitalLeague(id string) (*CapitalLeague, error) {
	data, err := h.do(http.MethodGet, CapitalLeaguesEndpoint.Build(id), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var league *CapitalLeague
	err = sonic.Unmarshal(data, &league)
	return league, err
}

// GetBuilderBaseLeague returns information about a single builder base league.
//
// GET /builderbaseleagues/{leagueId}
func (h *Client) GetBuilderBaseLeague(id string) (*BuilderBaseLeague, error) {
	data, err := h.do(http.MethodGet, BuilderBaseLeaguesEndpoint.Build(id), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var league *BuilderBaseLeague
	err = sonic.Unmarshal(data, &league)
	return league, err
}

// GetBuilderBaseLeagues returns a list of builder base leagues. Pass params=nil to get all leagues.
//
// GET /builderbaseleagues
func (h *Client) GetBuilderBaseLeagues(params *PagingParams) (*PaginatedResponse[BuilderBaseLeague], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, BuilderBaseLeaguesEndpoint.Build(), req, true)
	if err != nil {
		return nil, err
	}

	var leagues *PaginatedResponse[BuilderBaseLeague]
	err = sonic.Unmarshal(data, &leagues)
	return leagues, err
}

// GetLeague returns information about a single league.
//
// GET /leagues/{leagueId}
func (h *Client) GetLeague(id string) (*League, error) {
	data, err := h.do(http.MethodGet, LeaguesEndpoint.Build(id), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var league *League
	err = sonic.Unmarshal(data, &league)
	return league, err
}

// GetLeagueSeasons returns a list of league seasons. Pass params=nil to get all seasons.
//
// GET /leagues/{leagueId}/seasons
func (h *Client) GetLeagueSeasons(id int, params *PagingParams) (*PaginatedResponse[LeagueSeason], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, LeaguesEndpoint.Build(strconv.Itoa(id), "seasons"), req, true)
	if err != nil {
		return nil, err
	}

	var seasons *PaginatedResponse[LeagueSeason]
	err = sonic.Unmarshal(data, &seasons)
	return seasons, err
}

// GetWarLeague returns information about a single war league.
//
// GET /warleagues/{leagueId}
func (h *Client) GetWarLeague(id string) (*WarLeague, error) {
	data, err := h.do(http.MethodGet, WarLeaguesEndpoint.Build(id), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var league *WarLeague
	err = sonic.Unmarshal(data, &league)
	return league, err
}

// GetWarLeagues returns a list of war leagues. Pass params=nil to get all leagues.
//
// GET /warleagues
func (h *Client) GetWarLeagues(params *PagingParams) ([]*WarLeague, error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, WarLeaguesEndpoint.Build(), req, true)
	if err != nil {
		return nil, err
	}

	var leagues []*WarLeague
	err = sonic.Unmarshal(data, &leagues)
	return leagues, err
}
