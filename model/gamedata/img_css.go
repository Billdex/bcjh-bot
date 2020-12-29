package gamedata

type ObjImgInfo struct {
	Id     int
	X      int
	Y      int
	Width  int
	Height int
}

type ImgCSS struct {
	ChefImg   map[int]ObjImgInfo
	RecipeImg map[int]ObjImgInfo
	EquipImg  map[int]ObjImgInfo
}
