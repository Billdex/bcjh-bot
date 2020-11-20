package gamedata

type RecipeData struct {
	RecipeId  int    `json:"recipeId"`
	Name      string `json:"name"`
	GalleryId string `json:"galleryId"`
	Rarity    int    `json:"rarity"`
	Origin    string `json:"origin"`
	Stirfry   int    `json:"stirfry"`
	Bake      int    `json:"bake"`
	Boil      int    `json:"boil"`
	Steam     int    `json:"steam"`
	Fry       int    `json:"fry"`
	Cut       int    `json:"knife"`
	Price     int    `json:"price"`
	ExPrice   int    `json:"exPrice"`
	Gift      string `json:"gift"`
	Guests    []struct {
		Guest string `json:"guest"`
	} `json:"guests"`
	Limit     int    `json:"limit"`
	Time      int    `json:"time"`
	Unlock    string `json:"unlock"`
	Materials []struct {
		MaterialId int `json:"material"`
		Quantity   int `json:"quantity"`
	} `json:"materials"`
}
