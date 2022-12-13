package util

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

// LoadFont 载入字体，为字体文件路径时直接使用，否则视为路径获取目录下的唯一 ttf 文件
func LoadFont(path string) (*truetype.Font, error) {
	// 配置为字体文件时直接用，其他情况视为目录
	if !strings.HasSuffix(strings.ToLower(path), ".ttf") {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		fontNum := 0
		fileName := ""
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".ttf") {
				fontNum++
				fileName = file.Name()
			}
		}
		// 如果目录下仅有一个字体文件，直接载入，有多个时采用命名为 font.ttf 的字体文件
		if fontNum == 1 {
			path = fmt.Sprintf("%s/%s", path, fileName)
		} else {
			path = fmt.Sprintf("%s/%s", path, "font.ttf")
		}
	}

	return LoadFontFile(path)
}

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
