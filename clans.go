package clash

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/bytedance/sonic"
)

type Clan struct {
	WarLeague                   WarLeague     `json:"warLeague"`
	CapitalLeague               CapitalLeague `json:"capitalLeague"`
	MemberList                  []ClanMember  `json:"memberList"`
	Tag                         string        `json:"tag"`
	ChatLanguage                Language      `json:"chatLanguage"`
	BuilderBasePoints           int           `json:"clanBuilderBasePoints"`
	RequiredBuilderBaseTrophies int           `json:"requiredBuilderBaseTrophies"`
	RequiredTownHallLevel       int           `json:"requiredTownhallLevel"`
	IsFamilyFriendly            bool          `json:"IsFamilyFriendly"`
	IsWarLogPublic              bool          `json:"isWarLogPublic"`
	WarFrequency                string        `json:"warFrequency"`
	Level                       int           `json:"clanLevel"`
	WarWinStreak                int           `json:"warWinStreak"`
	WarWins                     int           `json:"warWins"`
	WarTies                     int           `json:"warTies"`
	WarLosses                   int           `json:"warLosses"`
	Points                      int           `json:"clanPoints"`
	CapitalPoints               int           `json:"clanCapitalPoints"`
	RequiredTrophies            int           `json:"requiredTrophies"`
	Labels                      []Label       `json:"labels"`
	Name                        string        `json:"name"`
	Location                    Location      `json:"location"`
	Type                        string        `json:"type"`
	MemberCount                 int           `json:"members"`
	Description                 string        `json:"description"`
	ClanCapital                 ClanCapital   `json:"clanCapital"`
	BadgeURLs                   ImageURLs     `json:"badgeUrls"`
}

type Clans []*Clan

func (c Clans) Tags() []string {
	tags := make([]string, len(c))
	for i, clan := range c {
		tags[i] = clan.Tag
	}
	return tags
}

type ClanMember struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

type ClanRole string

const (
	ClanRoleNotMember ClanRole = "notMember"
	ClanRoleMember    ClanRole = "member"
	ClanRoleAdmin     ClanRole = "admin"
	ClanRoleCoLeader  ClanRole = "coLeader"
	ClanRoleLeader    ClanRole = "leader"
)

func (r ClanRole) String() string {
	return string(r)
}

func (r ClanRole) Format() string {
	switch r {
	case ClanRoleNotMember:
		return "Not Member"
	case ClanRoleMember:
		return "Member"
	case ClanRoleAdmin:
		return "Admin"
	case ClanRoleCoLeader:
		return "Co-Leader"
	case ClanRoleLeader:
		return "Leader"
	default:
		return ""
	}
}

const (
	WarFrequencyUnknown       = "unknown"
	WarFrequencyAlways        = "always"
	WarFrequencyMTOncePerWeek = "moreThanOncePerWeek"
	WarFrequencyOncePerWeek   = "oncePerWeek"
	WarFrequencyLTOncePerWeek = "lessThanOncePerWeek"
	WarFrequencyNever         = "never"
	WarFrequencyAny           = "any"
)

type ClanCapital struct {
	CapitalHallLevel int            `json:"capitalHallLevel"`
	Districts        []ClanDistrict `json:"districts"`
}

type ClanDistrict struct {
	Name              string `json:"name"`
	ID                int    `json:"id"`
	DistrictHallLevel int    `json:"districtHallLevel"`
}

type ClanWar struct {
	Clan                 WarClan      `json:"clan"`
	Opponent             WarClan      `json:"opponent"`
	TeamSize             int          `json:"teamSize"`
	StartTime            string       `json:"startTime"`
	State                ClanWarState `json:"state"`
	EndTime              string       `json:"endTime"`
	PreparationStartTime string       `json:"preparationStartTime"`
}

type ClanWarState = string

const (
	ClanWarStateClanNotFound  ClanWarState = "clanNotFound"
	ClanWarStateAccessDenied  ClanWarState = "accessDenied"
	ClanWarStateNotInWar      ClanWarState = "notInWar"
	ClanWarStateInMatchmaking ClanWarState = "inMatchmaking"
	ClanWarStateEnterWar      ClanWarState = "enterWar"
	ClanWarStateMatched       ClanWarState = "matched"
	ClanWarStatePreparation   ClanWarState = "preparation"
	ClanWarStateWar           ClanWarState = "war"
	ClanWarStateInWar         ClanWarState = "inWar"
	ClanWarStateEnded         ClanWarState = "ended"
)

type ClanWarLeagueGroup struct {
	Tag    string                  `json:"tag"`
	State  ClanWarLeagueGroupState `json:"state"`
	Season string                  `json:"season"`
	Clans  []ClanWarLeagueClan
	Rounds []ClanWarLeagueRound
}

type ClanWarLeagueClan struct {
	Tag       string                    `json:"tag"`
	ClanLevel int                       `json:"clanLevel"`
	Name      string                    `json:"name"`
	Members   []ClanWarLeagueClanMember `json:"members"`
	BadgeURLs ImageURLs                 `json:"badgeUrls"`
}

type ClanWarLeagueClanMember struct {
	Tag           string `json:"tag"`
	TownHallLevel int    `json:"townHallLevel"`
	Name          string `json:"name"`
}

type ClanWarLeagueRound struct {
	WarTags []string `json:"warTags"`
}

type ClanWarLogEntry struct {
	Clan             WarClan       `json:"clan"`
	Opponent         WarClan       `json:"opponent"`
	TeamSize         int           `json:"teamSize"`
	AttacksPerMember int           `json:"attacksPerMember"`
	EndTime          string        `json:"endTime"`
	Result           ClanWarResult `json:"result"`
}

type WarClan struct {
	DestructionPercentage float64       `json:"destructionPercentage"`
	Tag                   string        `json:"tag"`
	Name                  string        `json:"name"`
	BadgeURLs             ImageURLs     `json:"badgeUrls"`
	ClanLevel             int           `json:"clanLevel"`
	Attacks               int           `json:"attacks"`
	Stars                 int           `json:"stars"`
	ExpEarned             int           `json:"expEarned"`
	Members               ClanWarMember `json:"members"`
}

type ClanWarMember struct {
	Tag                string          `json:"tag"`
	Name               string          `json:"name"`
	MapPosition        int             `json:"mapPosition"`
	TownHallLevel      int             `json:"townHallLevel"`
	OpponentAttacks    int             `json:"opponentAttacks"`
	BestOpponentAttack *ClanWarAttack  `json:"bestOpponentAttack,omitempty"`
	Attacks            []ClanWarAttack `json:"attacks,omitempty"`
}

type ClanWarAttack struct {
	Order                 int    `json:"order"`
	AttackerTag           string `json:"attackerTag"`
	DefenderTag           string `json:"defenderTag"`
	Stars                 int    `json:"stars"`
	DestructionPercentage int    `json:"destructionPercentage"`
	Duration              int    `json:"duration"`
}

type SearchClanParams struct {
	*PagingParams
	Name          string   `json:"name,omitempty"`
	WarFrequency  string   `json:"warFrequency,omitempty"`
	LocationID    string   `json:"locationId,omitempty"`
	MinMembers    string   `json:"minMembers,omitempty"`
	MaxMembers    string   `json:"maxMembers,omitempty"`
	MinClanPoints string   `json:"minClanPoints,omitempty"`
	MinClanLevel  string   `json:"minClanLevel,omitempty"`
	LabelIDs      []string `json:"labelIds,omitempty"`
}

type ClanWarLeagueGroupState = string

const (
	ClanWarLeagueGroupStateNotFound ClanWarLeagueGroupState = "groupNotFound"
	ClanWarLeagueGroupStateNotInWar ClanWarLeagueGroupState = "notInWar"
	ClanWarLeagueGroupStatePrep     ClanWarLeagueGroupState = "preparation"
	ClanWarLeagueGroupStateWar      ClanWarLeagueGroupState = "war"
	ClanWarLeagueGroupStateEnded    ClanWarLeagueGroupState = "ended"
)

type ClanWarResult = string

const (
	ClanWarResultWin  ClanWarResult = "win"
	ClanWarResultLose ClanWarResult = "lose"
	ClanWarResultTie  ClanWarResult = "tie"
)

type ClanCapitalRaidSeason struct {
	AttackLog               []ClanCapitalRaidSeasonAttackLogEntry  `json:"attackLog"`
	DefenseLog              []ClanCapitalRaidSeasonDefenseLogEntry `json:"defenseLog"`
	State                   string                                 `json:"state"`
	StartTime               string                                 `json:"startTime"`
	EndTime                 string                                 `json:"endTime"`
	CapitalTotalLoot        int                                    `json:"capitalTotalLoot"`
	RaidsCompleted          int                                    `json:"raidsCompleted"`
	TotalAttacks            int                                    `json:"totalAttacks"`
	EnemyDistrictsDestroyed int                                    `json:"enemyDistrictsDestroyed"`
	OffensiveReward         int                                    `json:"offensiveReward"`
	DefensiveReward         int                                    `json:"defensiveReward"`
	Members                 []ClanCapitalRaidSeasonMember          `json:"members"`
}

type ClanCapitalRaidSeasonAttackLogEntry struct {
	Defender           ClanCapitalRaidSeasonClanInfo   `json:"defender"`
	AttackCount        int                             `json:"attackCount"`
	DistrictCount      int                             `json:"districtCount"`
	DistrictsDestroyed int                             `json:"districtsDestroyed"`
	Districts          []ClanCapitalRaidSeasonDistrict `json:"districts"`
}

type ClanCapitalRaidSeasonDefenseLogEntry struct {
	Attacker           ClanCapitalRaidSeasonClanInfo   `json:"attacker"`
	AttackCount        int                             `json:"attackCount"`
	DistrictCount      int                             `json:"districtCount"`
	DistrictsDestroyed int                             `json:"districtsDestroyed"`
	Districts          []ClanCapitalRaidSeasonDistrict `json:"districts"`
}

type ClanCapitalRaidSeasonMember struct {
	Tag                    string `json:"tag"`
	Name                   string `json:"name"`
	Attacks                int    `json:"attacks"`
	AttackLimit            int    `json:"attackLimit"`
	BonusAttackLimit       int    `json:"bonusAttackLimit"`
	CapitalResourcesLooted int    `json:"capitalResourcesLooted"`
}

type ClanCapitalRaidSeasonClanInfo struct {
	Tag       string    `json:"tag"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	BadgeURLs ImageURLs `json:"badgeUrls"`
}

type ClanCapitalRaidSeasonDistrict struct {
	Stars              int                           `json:"stars"`
	Name               string                        `json:"name"`
	ID                 int                           `json:"id"`
	DestructionPercent int                           `json:"destructionPercent"`
	AttackCount        int                           `json:"attackCount"`
	TotalLooted        int                           `json:"totalLooted"`
	Attacks            []ClanCapitalRaidSeasonAttack `json:"attacks"`
	DistrictHallLevel  int                           `json:"districtHallLevel"`
}

type ClanCapitalRaidSeasonAttack struct {
	Attacker           ClanCapitalRaidSeasonAttacker `json:"attacker"`
	DestructionPercent int                           `json:"destructionPercent"`
	Stars              int                           `json:"stars"`
}

type ClanCapitalRaidSeasonAttacker struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

// GetCurrentClanWarLeagueGroup returns the current war league group for a clan.
//
// GET /clans/{clanTag}/currentwar/leaguegroup
func (h *Client) GetCurrentClanWarLeagueGroup(tag string) (*ClanWarLeagueGroup, error) {
	tag = TagURLSafe(CorrectTag(tag))
	data, err := h.do(http.MethodGet, ClansEndpoint.Build(tag, "currentwar/leaguegroup"), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var group *ClanWarLeagueGroup
	err = sonic.Unmarshal(data, &group)
	return group, err
}

// GetClanWarLeagueWar returns information about a single war within a clan war league.
//
// GET /clanwarleagues/wars/{warTag}
func (h *Client) GetClanWarLeagueWar(warTag string) (*ClanWarLeagueGroup, error) {
	data, err := h.do(http.MethodGet, ClanWarLeaguesEndpoint.Build("wars", warTag), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var war *ClanWarLeagueGroup
	err = sonic.Unmarshal(data, &war)
	return war, err
}

// GetClanWarLog returns a clan's war log.
//
// GET /clans/{clanTag}/warlog
func (h *Client) GetClanWarLog(tag string, params *PagingParams) (*PaginatedResponse[ClanWarLogEntry], error) {
	tag = TagURLSafe(CorrectTag(tag))
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, ClansEndpoint.Build(tag, "warlog"), req, true)
	if err != nil {
		return nil, err
	}

	var log *PaginatedResponse[ClanWarLogEntry]
	err = sonic.Unmarshal(data, &log)
	return log, err
}

// SearchClans returns a list of clans that match the given params.
//
// GET /clans
func (h *Client) SearchClans(params SearchClanParams) (*PaginatedResponse[Clan], error) {
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params.PagingParams).
		SetQueryParamsFromValues(url.Values{
			"name":          {params.Name},
			"warFrequency":  {params.WarFrequency},
			"locationId":    {params.LocationID},
			"minMembers":    {params.MinMembers},
			"maxMembers":    {params.MaxMembers},
			"minClanPoints": {params.MinClanPoints},
			"minClanLevel":  {params.MinClanLevel},
			"labelIds":      params.LabelIDs,
		})
	data, err := h.do(http.MethodGet, ClansEndpoint.Build(), req, true)
	if err != nil {
		return nil, err
	}

	var clans *PaginatedResponse[Clan]
	err = sonic.Unmarshal(data, &clans)
	return clans, err
}

// GetCurrentClanWar returns information about a clan's current clan war.
//
// GET /clans/{clanTag}/currentwar
func (h *Client) GetCurrentClanWar(tag string) (*ClanWar, error) {
	tag = TagURLSafe(CorrectTag(tag))
	data, err := h.do(http.MethodGet, ClansEndpoint.Build(tag, "currentwar"), h.withAuth(h.newDefaultRequest()), true)
	if err != nil {
		return nil, err
	}

	var war *ClanWar
	err = sonic.Unmarshal(data, &war)
	return war, err
}

// GetClan returns a clan by its tag.
//
// GET /clans/{clanTag}
func (h *Client) GetClan(tag string) (*Clan, error) {
	tag = TagURLSafe(CorrectTag(tag))
	req := h.withAuth(h.newDefaultRequest())
	data, err := h.do(http.MethodGet, ClansEndpoint.Build(tag), req, true)
	if err != nil {
		return nil, err
	}

	var clan *Clan
	err = sonic.Unmarshal(data, &clan)
	return clan, err
}

// GetClans makes use of concurrency to get multiple clans simultaneously. The original order of the tags is preserved.
func (h *Client) GetClans(tags ...string) (Clans, error) {
	var wg sync.WaitGroup
	clans := make(Clans, len(tags))
	errChan := make(chan error, len(tags))

	for i, tag := range tags {
		wg.Add(1)
		go func(i int, tag string) {
			defer wg.Done()
			clan, err := h.GetClan(tag)
			if err != nil {
				errChan <- err
				return
			}
			clans[i] = clan
		}(i, tag)
	}
	wg.Wait()

	if len(errChan) > 0 {
		return nil, <-errChan
	}
	return clans, nil
}

func (h *Client) GetClanMembers(tag string, params *PagingParams) (*PaginatedResponse[ClanMember], error) {
	tag = TagURLSafe(CorrectTag(tag))
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, ClansEndpoint.Build(tag, "members"), req, true)
	if err != nil {
		return nil, err
	}

	var members *PaginatedResponse[ClanMember]
	err = sonic.Unmarshal(data, &members)
	return members, err
}

func (h *Client) GetClanCapitalRaidSeasons(tag string, params *PagingParams) (*PaginatedResponse[ClanCapitalRaidSeason], error) {
	tag = TagURLSafe(CorrectTag(tag))
	req := h.withPaging(h.withAuth(h.newDefaultRequest()), params)
	data, err := h.do(http.MethodGet, ClansEndpoint.Build(tag, "capitalraidseasons"), req, true)
	if err != nil {
		return nil, err
	}

	var seasons *PaginatedResponse[ClanCapitalRaidSeason]
	err = sonic.Unmarshal(data, &seasons)
	return seasons, err
}
