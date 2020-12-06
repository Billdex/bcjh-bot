package gamedata

type EquipData struct {
	EquipId   int    `json:"equipId"`
	Name      string `json:"name"`
	GalleryId string `json:"galleryId"`
	Origin    string `json:"origin"`
	Rarity    int    `json:"rarity"`
	Skills    []int  `json:"skill"`
}
