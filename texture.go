package main

import (
	"fmt"
	"image"
	"os"
	"image/png"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	TEX_DEFAULT_INTERAL_FORMAT = gl.RGBA
	TEX_DEFAULT_IMAGE_FORMAT = gl.RGBA
	TEX_DEFAULT_WRAP_S = gl.REPEAT
	TEX_DEFAULT_WRAP_T = gl.REPEAT
	TEX_DEFAULT_FILTER_MIN = gl.LINEAR
	TEX_DEFAULT_FILTER_MAX = gl.LINEAR
)

type Texture struct {
	ID uint32
	Width, Height int32
	InternalFormat int32
	ImageFormat uint32

	WrapS, WrapT int32
	FilterMin, FilterMax int32
}

func NewTexture(texPath string, shader *shader) (*Texture, error) {
	newTexture := new(Texture)
	gl.GenTextures(1, &newTexture.ID)
	newTexture.SetDefaults()
	
	width, height, img, err := LoadTextureImg(texPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to create texture: %v", err)
	}

	newTexture.Width, newTexture.Height = width, height
	newTexture.Generate(img)

	shader.use()
	shader.SetInt("texture1", 0)
	return newTexture, nil
}


func LoadTextureImg(texPath string) (int32, int32, []uint8, error) {
	imgFile, err := os.Open(texPath)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to decode image: %v", err)
	}
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	
	//extract the data
	// pixelData := ExtractImgData(img, width, height)
	switch convImg := img.(type) {
	case *image.RGBA:
		return int32(width), int32(height), convImg.Pix, nil
	case *image.NRGBA:
		return int32(width), int32(height), convImg.Pix, nil
	}
	//flipImage(img.(*image.RGBA))
	return int32(width), int32(height), img.(*image.RGBA).Pix, nil
	
}

// SetDefaults sets the teextures fields to be the defualt values
func (t *Texture) SetDefaults()  {
	t.InternalFormat = TEX_DEFAULT_INTERAL_FORMAT
	t.ImageFormat = TEX_DEFAULT_IMAGE_FORMAT
	t.WrapS = TEX_DEFAULT_WRAP_S
	t.WrapT = TEX_DEFAULT_WRAP_T
	t.FilterMin = TEX_DEFAULT_FILTER_MIN
	t.FilterMax = TEX_DEFAULT_FILTER_MAX
}

// Generate generates the texture data
func (t *Texture) Generate(data []uint8)  {
	// make the texture
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, t.InternalFormat, t.Width, t.Height, 0, t.ImageFormat, gl.UNSIGNED_BYTE, gl.Ptr(data))

	//set the texture wrap and filter stuff
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, t.WrapS)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, t.WrapT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, t.FilterMin)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, t.FilterMax)

	//unbind the texture
	// gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Bind binds the texture
func (t *Texture) Bind(texUnit uint32)  {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

func flipImage(img *image.RGBA) {
	width, height := img.Bounds().Size().X, img.Bounds().Size().Y
	for y := 0; y < height/2; y++ {
		for x := 0; x < width; x++ {
			tmp := img.RGBAAt(x, y)
			img.SetRGBA(x, y, img.RGBAAt(x, height-1-y))
			img.SetRGBA(x, height-1-y, tmp)
		}
	}
}
