package image

import (
	"bufio"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

type InkMat struct {
	W int
	H int
	P []byte
}

func LoadImage(path string) (*InkMat, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}

	mat := new(InkMat)
	mat.W = img.Bounds().Dx()
	mat.H = img.Bounds().Dy()
	mat.P = make([]byte, mat.W*mat.H)

	// グレイ化
	idx := 0
	for i := 0; i < mat.H; i++ {
		for j := 0; j < mat.W; j++ {
			gray, _, _, alpha := color.GrayModel.Convert(img.At(j, i)).RGBA()

			if alpha == 0 {
				mat.P[idx] = 0
			} else {
				outa := math.Min(float64(alpha)/0xffff, 1.0)
				outc := math.Min(float64(gray)+(1-outa), 1.0)

				mat.P[idx] = byte(outc * 0xff)
			}
			idx++
		}
	}

	return mat, nil
}
