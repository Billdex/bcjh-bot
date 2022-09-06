package util

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

// LoadFontFile 载入字体文件
func LoadFontFile(path string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return font, nil
}

// LoadPngImageFile 载入png格式的图片文件
func LoadPngImageFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(f)
	return img, err
}

// SavePngImage 将图片按照 png 格式保存
func SavePngImage(path string, img image.Image) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	err = png.Encode(dst, img)
	if err != nil {
		return err
	}
	_ = dst.Close()
	return nil
}
