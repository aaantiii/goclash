package clash

import (
	"net/http"
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

// GetClan returns a clan by its tag.
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

// GetClans makes use of concurrency to get multiple clans simultaneously.
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
