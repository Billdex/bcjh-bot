package util

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/draw"
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

// LoadImageFile 载入图片文件
func LoadImageFile(path string, width int, height int, op draw.Op) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fImg, _ := png.Decode(f)
	_ = f.Close()

	draw.Draw(img, img.Bounds(), fImg, fImg.Bounds().Min, op)
	return img, nil
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
