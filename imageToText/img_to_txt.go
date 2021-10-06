package imageToText

import (
	"bytes"
	"github.com/nfnt/resize"
	"image/png"
	"os"
)

var(
	DefaultTxt=[]string{"@", "#", "*", "%", "+","~", ",", ".", " "}
)
func ImgToTxt(imgPath string, width,height int,rowEnd string, txt []string)(string,error){
	file, err := os.Open(imgPath)
	if err != nil {
		return "",err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return "",err
	}
	newImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	dx := newImg.Bounds().Dx()
	dy := newImg.Bounds().Dy()
	textBuffer := bytes.Buffer{}
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			colorRgb := newImg.At(x, y)
			r, g, b, a := colorRgb.RGBA()
			v := int(float64(r+g+b)/3*float64(a>>8)/255.0)>>8
			num := int(float64(v)/256*float64(len(txt)))
			textBuffer.WriteString(txt[num])
		}
		textBuffer.WriteString(rowEnd)
	}
	return textBuffer.String(),nil
}