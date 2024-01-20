package clash

type WarLeague struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type WarLeagueID int

const (
	WarLeagueUnranked WarLeagueID = 48000000 + iota
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

type CapitalLeague struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type LeagueData struct {
	Paging  Paging   `json:"paging,omitempty"`
	Leagues []League `json:"items,omitempty"`
}

type League struct {
	IconUrls ImageURLs `json:"iconUrls,omitempty"`
	Name     string    `json:"name,omitempty"`
	ID       int       `json:"id,omitempty"`
}
