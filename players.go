package goclash

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
)

// PlayerBase is embedded in Player and contains the most basic information about a player. May be used as DTO.
type PlayerBase struct {
	League                   League            `json:"league"`
	BuilderBaseLeague        BuilderBaseLeague `json:"builderBaseLeague"`
	Clan                     PlayerClan        `json:"clan"`
	Role                     ClanRole          `json:"role"`
	AttackWins               int               `json:"attackWins"`
	DefenseWins              int               `json:"defenseWins"`
	TownHallLevel            int               `json:"townHallLevel"`
	Tag                      string            `json:"tag"`
	Name                     string            `json:"name"`
	ExpLevel                 int               `json:"expLevel"`
	Trophies                 int               `json:"trophies"`
	BestTrophies             int               `json:"bestTrophies"`
	Donations                int               `json:"donations"`
	DonationsReceived        int               `json:"donationsReceived"`
	BuilderHallLevel         int               `json:"builderHallLevel"`
	BuilderBaseTrophies      int               `json:"builderBaseTrophies"`
	BestBuilderBaseTrophies  int               `json:"bestBuilderBaseTrophies"`
	WarStars                 int               `json:"warStars"`
	ClanCapitalContributions int               `json:"clanCapitalContributions"`
}

// Player is a player returned by the API.
type Player struct {
	*PlayerBase
	WarPreference       string            `json:"warPreference"`
	TownHallWeaponLevel int               `json:"townHallWeaponLevel"`
	LegendStatistics    LegendStatistics  `json:"legendStatistics"`
	Troops              []PlayerItemLevel `json:"troops"`
	Heroes              []PlayerItemLevel `json:"heroes,omitempty"`
	HeroEquipment       []PlayerItemLevel `json:"heroEquipment,omitempty"`
	Spells              []PlayerItemLevel `json:"spells"`
	Labels              []Label           `json:"labels"`
	Achievements        []Achievement     `json:"achievements"`
	PlayerHouse         PlayerHouse       `json:"playerHouse"`
}

// InGameURL returns a link.clashofclans.com URL that can be used to open the player profile in game.
func (p *Player) InGameURL() string {
	return "https://link.clashofclans.com?action=OpenPlayerProfile&tag=" + TagURLSafe(p.Tag)
}

// GetAchievement returns an IndexedAchievement by Achievement.Name and Achievement.Info. The index can be used to get the same achievement from other players, to make it more efficient.
func (p *Player) GetAchievement(achievement *Achievement) (*IndexedAchievement, error) {
	for i, a := range p.Achievements {
		switch a.Name {
		case AchievementKeepYourAccountSafeOld.Name:
			if a.Name != achievement.Name || a.Info != achievement.Info {
				continue
			}
		default:
			if a.Name != achievement.Name {
				continue
			}
		}
		return &IndexedAchievement{
			Achievement: &a,
			Index:       i,
		}, nil
	}
	return nil, errors.New("achievement not found")
}

type Players []*Player

func (p Players) Tags() []string {
	tags := make([]string, len(p))
	for i, player := range p {
		tags[i] = player.Tag
	}
	return tags
}

// String implements fmt.Stringer, returning a comma-separated list of player names.
func (p Players) String() string {
	str := make([]string, len(p))
	for i, player := range p {
		str[i] = player.Name
	}
	return strings.Join(str, ", ")
}

func (p Players) GetAchievement(achievement *Achievement) ([]*Achievement, error) {
	if len(p) == 0 {
		return nil, errors.New("no players were provided")
	}
	indexed, err := p[0].GetAchievement(achievement)
	if err != nil {
		return nil, err
	}

	achievements := make([]*Achievement, len(p))
	for i, player := range p {
		achievements[i] = &player.Achievements[indexed.Index]
	}
	return achievements, nil
}

type PlayerClan struct {
	Tag       string    `json:"tag"`
	Level     int       `json:"clanLevel"`
	Name      string    `json:"name"`
	BadgeURLs ImageURLs `json:"badgeUrls"`
}

type PlayerItemLevel struct {
	Name               string            `json:"name"`
	Village            string            `json:"village"`
	Level              int               `json:"level"`
	MaxLevel           int               `json:"maxLevel"`
	SuperTroopIsActive bool              `json:"superTroopIsActive,omitempty"`
	Equipment          []PlayerItemLevel `json:"equipment,omitempty"`
}

type LegendStatistics struct {
	PreviousSeason   Season        `json:"previousSeason,omitempty"`
	BestSeason       Season        `json:"bestSeason,omitempty"`
	BestVersusSeason Season        `json:"bestVersusSeason,omitempty"`
	CurrentSeason    CurrentSeason `json:"currentSeason,omitempty"`
	LegendTrophies   int           `json:"legendTrophies,omitempty"`
}

type CurrentSeason struct {
	Rank     int `json:"rank,omitempty"`
	Trophies int `json:"trophies,omitempty"`
}

type Season struct {
	ID       string `json:"id,omitempty"`
	Rank     int    `json:"rank,omitempty"`
	Trophies int    `json:"trophies,omitempty"`
}

type BuilderBaseLeague struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerVerification struct {
	Tag    string `json:"tag"`
	Token  string `json:"token"`
	Status string `json:"status"`
}

type PlayerHouse struct {
	Elements []PlayerHouseElement `json:"elements"`
}

type PlayerHouseElement struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

func (v *PlayerVerification) IsOk() bool {
	return v.Status == PlayerVerificationStatusOk
}

const (
	VillageHome    = "home"
	VillageBuilder = "builderBase"

	PlayerVerificationStatusOk      = "ok"
	PlayerVerificationStatusInvalid = "invalid"

	PlayerHouseElementTypeGround = "ground"
	PlayerHouseElementTypeRoof   = "roof"
	PlayerHouseElementTypeFoot   = "foot"
	PlayerHouseElementTypeDeco   = "deco"
)

// GetPlayer returns information about a single player by tag.
//
// GET /players/{playerTag}
func (h *Client) GetPlayer(tag string) (*Player, error) {
	tag = TagURLSafe(CorrectTag(tag))
	req := h.withAuth(h.newDefaultRequest())
	data, err := h.do(http.MethodGet, PlayersEndpoint.Build(tag), req, true)
	if err != nil {
		return nil, err
	}

	var player *Player
	err = sonic.Unmarshal(data, &player)
	return player, err
}

// GetPlayers makes use of concurrency to get multiple players simultaneously. Players that failed to be fetched will be nil in the returned slice.
func (h *Client) GetPlayers(tags ...string) Players {
	var wg sync.WaitGroup
	players := make(Players, len(tags))
	for i, tag := range tags {
		wg.Add(1)
		go func(i int, tag string) {
			defer wg.Done()
			player, _ := h.GetPlayer(tag)
			players[i] = player
		}(i, tag)
	}
	wg.Wait()
	return players
}

// GetPlayersWithError makes use of concurrency to get multiple players simultaneously. Unlike GetPlayers, this function returns nil and error if any of the players failed to be fetched.
func (h *Client) GetPlayersWithError(tags ...string) (Players, error) {
	var wg sync.WaitGroup
	players := make(Players, len(tags))
	errChan := make(chan error, len(tags))

	for i, tag := range tags {
		wg.Add(1)
		go func(i int, tag string) {
			defer wg.Done()
			player, err := h.GetPlayer(tag)
			if err != nil {
				errChan <- err
				return
			}
			players[i] = player
		}(i, tag)
	}
	wg.Wait()

	if len(errChan) > 0 {
		return nil, <-errChan
	}
	return players, nil
}

// VerifyPlayer verifies a player token.
//
// POST /players/{playerTag}/verifytoken
func (h *Client) VerifyPlayer(tag, token string) (*PlayerVerification, error) {
	tag = TagURLSafe(CorrectTag(tag))
	req := h.withAuth(h.newDefaultRequest()).SetBody(map[string]string{
		"token": token,
	})
	data, err := h.do(http.MethodPost, PlayersEndpoint.Build(fmt.Sprintf("%s/verifytoken", tag)), req, false)
	if err != nil {
		return nil, err
	}

	var verification *PlayerVerification
	err = sonic.Unmarshal(data, &verification)
	return verification, err
}
