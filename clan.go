package clash

type Clan struct {
	Tag        string        `json:"tag"`
	Name       string        `json:"name"`
	MemberList []*ClanMember `json:"memberList"`
	Members    int           `json:"members"`
}

type ClanMember struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}
