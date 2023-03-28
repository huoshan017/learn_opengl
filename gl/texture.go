package gl

import (
	"log"
	"unsafe"

	"github.com/huoshan017/go-stbi"
)

func TextureFromFile(path, directory string, gamma bool) uint32 {
	var filename = path
	filename = directory + "/" + filename

	var textureId uint32
	GenTextures(1, &textureId)

	errCode := GetError()
	if errCode != 0 {
		log.Fatalf("fatal!!! get error %v", errCode)
	}

	if textureId == 0 {
		log.Fatalf("fatal!!! generate texture id zero")
	}

	var nChannels int32
	image, err := stbi.Load(filename, &nChannels, 0)
	if err != nil {
		log.Fatalf("Texture failed to load at path %v with directory %v, err %v", path, directory, err)
	}

	var format int32
	if nChannels == 1 {
		format = RED
	} else if nChannels == 3 {
		format = RGB
	} else if nChannels == 4 {
		format = RGBA
	}

	BindTexture(TEXTURE_2D, textureId)
	width := image.Rect.Dx()
	height := image.Rect.Dy()
	TexImage2D(TEXTURE_2D, 0, format, int32(width), int32(height), 0, uint32(format), UNSIGNED_BYTE, unsafe.Pointer(&image.Pix[0]))
	GenerateMipmap(TEXTURE_2D)

	TexParameteri(TEXTURE_2D, TEXTURE_WRAP_S, REPEAT)
	TexParameteri(TEXTURE_2D, TEXTURE_WRAP_T, REPEAT)
	TexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, LINEAR_MIPMAP_LINEAR)
	TexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, LINEAR)

	return textureId
}
