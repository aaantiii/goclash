package clash

type Player struct {
	Tag          string         `json:"tag"`
	Name         string         `json:"name"`
	Achievements []*Achievement `json:"achievements"`
	Clan         *Clan          `json:"clan"`
}
