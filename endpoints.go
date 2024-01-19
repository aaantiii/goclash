package clash

type Endpoint string

func (e Endpoint) URL() string {
	return BaseURL + string(e)
}

const (
	BaseURL                     = "https://api.clashofclans.com/v1"
	ClansEndpoint      Endpoint = "/clans"
	PlayersEndpoint    Endpoint = "/players"
	LeaguesEndpoint    Endpoint = "/leagues"
	WarLeaguesEndpoint Endpoint = "/warleagues"
	LocationsEndpoint  Endpoint = "/locations"
	GoldpassEndpoint   Endpoint = "/goldpass/seasons/current"
	LabelsEndpoint     Endpoint = "/labels"
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
