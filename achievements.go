package clash

// Achievement represents a Clash of Clans achievement.
// Use AchievementIndex* constants to index into the Achievements slice.
type Achievement struct {
	Name           string `json:"name"`
	Stars          int    `json:"stars"`
	Value          int    `json:"value"`
	Target         int    `json:"target"`
	Info           string `json:"info"`
	CompletionInfo string `json:"completionInfo"`
	Village        string `json:"village"`
}

const (
	AchievementIndexBiggerCoffers = iota
	AchievementIndexGetThoseGoblins
	AchievementIndexBiggerAndBetter
	AchievementIndexNiceAndTidy
	AchievementIndexDiscoverNewTroops
	AchievementIndexGoldGrab
	AchievementIndexElixirEscapade
	AchievementIndexSweetVictory
	AchievementIndexEmpireBuilder
	AchievementIndexWallBuster
	AchievementIndexHumiliator
	AchievementIndexUnionBuster
	AchievementIndexConqueror
	AchievementIndexUnbreakable
	AchievementIndexFriendInNeed
	AchievementIndexMortarMauler
	AchievementIndexHeroicHeist
	AchievementIndexLeagueAllStar
	AchievementIndexXbowExterminator
	AchievementIndexFirefighter
	AchievementIndexWarHero
	AchievementIndexClanWarWealth
	AchievementIndexAntiArtillery
	AchievementIndexSharingIsCaring
	AchievementIndexKeepYourAccountSafeOld
	AchievementIndexMasterEngineering
	AchievementIndexNextGenerationModel
	AchievementIndexUnBuildIt
	AchievementIndexChampionBuilder
	AchievementIndexHighGear
	AchievementIndexHiddenTreasures
	AchievementIndexGamesChampion
	AchievementIndexDragonSlayer
	AchievementIndexWarLeagueLegend
	AchievementIndexKeepYourAccountSafe
	AchievementIndexWellSeasoned
	AchievementIndexShatteredAndScattered
	AchievementIndexNotSoEasyThisTime
	AchievementIndexBustThis
	AchievementIndexSuperbWork
	AchievementIndexSiegeSharer
	AchievementIndexAggressiveCapitalism
	AchievementIndexMostValuableClanmate
	AchievementIndexCounterspell
	AchievementIndexMonolithMasher
	AchievementIndexUngratefulChild
)
