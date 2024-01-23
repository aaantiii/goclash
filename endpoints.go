package goclash

type Endpoint string

// Build returns the full URL for the endpoint.
//
// Example: PlayersEndpoint.Build("ABC123") returns "https://api.clashofclans.com/v1/players/ABC123"
func (e Endpoint) Build(routes ...string) string {
	url := BaseURL + string(e)
	for _, route := range routes {
		url += "/" + route
	}

	return url
}

const (
	BaseURL                             = "https://api.clashofclans.com/v1"
	ClansEndpoint              Endpoint = "/clans"
	ClanWarLeaguesEndpoint     Endpoint = "/clanwarleagues"
	PlayersEndpoint            Endpoint = "/players"
	LeaguesEndpoint            Endpoint = "/leagues"
	WarLeaguesEndpoint         Endpoint = "/warleagues"
	BuilderBaseLeaguesEndpoint Endpoint = "/builderbaseleagues"
	CapitalLeaguesEndpoint     Endpoint = "/capitalleagues"
	LocationsEndpoint          Endpoint = "/locations"
	GoldPassEndpoint           Endpoint = "/goldpass/seasons/current"
	LabelsEndpoint             Endpoint = "/labels"
)

type DevEndpoint string

func (e DevEndpoint) URL() string {
	return DevBaseURL + string(e)
}

const (
	DevBaseURL                       = "https://developer.clashofclans.com"
	DevLoginEndpoint     DevEndpoint = "/api/login"
	DevKeyListEndpoint   DevEndpoint = "/api/apikey/list"
	DevKeyCreateEndpoint DevEndpoint = "/api/apikey/create"
	DevKeyRevokeEndpoint DevEndpoint = "/api/apikey/revoke"
	IPifyEndpoint                    = "https://api.ipify.org"
)
