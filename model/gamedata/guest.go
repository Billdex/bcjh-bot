package gamedata

type GuestData struct {
	Name  string `json:"name"`
	Gifts []struct {
		Antique string `json:"antique"`
		Recipe  string `json:"recipe"`
	} `json:"gifts"`
}
