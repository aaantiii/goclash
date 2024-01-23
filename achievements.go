package goclash

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

// IndexedAchievement embeds Achievement and adds the index of the achievement in the Player.Achievements slice to it.
type IndexedAchievement struct {
	*Achievement
	Index int
}

var (
	AchievementBiggerCoffers           = &Achievement{Name: "Bigger Coffers"}
	AchievementGetThoseGoblins         = &Achievement{Name: "Get those Goblins!"}
	AchievementBiggerAndBetter         = &Achievement{Name: "Bigger & Better"}
	AchievementNiceAndTidy             = &Achievement{Name: "Nice and Tidy"}
	AchievementDiscoverNewTroops       = &Achievement{Name: "Discover New Troops"}
	AchievementGoldGrab                = &Achievement{Name: "Gold Grab"}
	AchievementElixirEscapade          = &Achievement{Name: "Elixir Escapade"}
	AchievementSweetVictory            = &Achievement{Name: "Sweet Victory!"}
	AchievementEmpireBuilder           = &Achievement{Name: "Empire Builder"}
	AchievementWallBuster              = &Achievement{Name: "Wall Buster"}
	AchievementHumiliator              = &Achievement{Name: "Humiliator"}
	AchievementUnionBuster             = &Achievement{Name: "Union Buster"}
	AchievementConqueror               = &Achievement{Name: "Conqueror"}
	AchievementUnbreakable             = &Achievement{Name: "Unbreakable"}
	AchievementFriendInNeed            = &Achievement{Name: "Friend in Need"}
	AchievementMortarMauler            = &Achievement{Name: "Mortar Mauler"}
	AchievementHeroicHeist             = &Achievement{Name: "Heroic Heist"}
	AchievementLeagueAllStar           = &Achievement{Name: "League All-Star"}
	AchievementXBowExterminator        = &Achievement{Name: "X-Bow Exterminator"}
	AchievementFirefighter             = &Achievement{Name: "Firefighter"}
	AchievementWarHero                 = &Achievement{Name: "War Hero"}
	AchievementClanWarWealth           = &Achievement{Name: "Clan War Wealth"}
	AchievementAntiArtillery           = &Achievement{Name: "Anti-Artillery"}
	AchievementSharingIsCaring         = &Achievement{Name: "Sharing is caring"}
	AchievementKeepYourAccountSafeOld  = &Achievement{Name: "Keep Your Account Safe!", Info: "Protect your village by connecting to a social network"}
	AchievementMasterEngineering       = &Achievement{Name: "Master Engineering"}
	AchievementNextGenerationModel     = &Achievement{Name: "Next Generation Model"}
	AchievementUnBuildIt               = &Achievement{Name: "Un-Build It"}
	AchievementChampionBuilder         = &Achievement{Name: "Champion Builder"}
	AchievementHighGear                = &Achievement{Name: "High Gear"}
	AchievementHiddenTreasures         = &Achievement{Name: "Hidden Treasures"}
	AchievementGamesChampion           = &Achievement{Name: "Games Champion"}
	AchievementDragonSlayer            = &Achievement{Name: "Dragon Slayer"}
	AchievementWarLeagueLegend         = &Achievement{Name: "War League Legend"}
	AchievementKeepYourAccountSafeSCID = &Achievement{Name: "Keep Your Account Safe!", Info: "Connect your account to Supercell ID for safe keeping."}
	AchievementWellSeasoned            = &Achievement{Name: "Well Seasoned"}
	AchievementShatteredAndScattered   = &Achievement{Name: "Shattered and Scattered"}
	AchievementNotSoEasyThisTime       = &Achievement{Name: "Not So Easy This Time"}
	AchievementBustThis                = &Achievement{Name: "Bust This"}
	AchievementSuperbWork              = &Achievement{Name: "Superb Work"}
	AchievementSiegeSharer             = &Achievement{Name: "Siege Sharer"}
	AchievementCounterspell            = &Achievement{Name: "Counterspell"}
	AchievementMonolithMasher          = &Achievement{Name: "Monolith Masher"}
	AchievementGetThoseOtherGoblins    = &Achievement{Name: "Get those other Goblins!"}
	AchievementGetEvenMoreGoblins      = &Achievement{Name: "Get even more Goblins!"}
	AchievementUngratefulChild         = &Achievement{Name: "Ungrateful Child"}
	AchievementAggressiveCapitalism    = &Achievement{Name: "Aggressive Capitalism"}
	AchievementMostValuableClanmate    = &Achievement{Name: "Most Valuable Clanmate"}
)
