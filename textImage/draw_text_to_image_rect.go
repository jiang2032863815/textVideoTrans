package textImage

import (
	"github.com/golang/freetype"
	"image"
	"image/color"
)

func DrawTextToImageRect(text string,textXCount,textYCount,imageHeight int,fontBytes []byte)(*image.RGBA,error){
	var alpha=0.65
	var textHeight = float64(imageHeight)/float64(textYCount)
	var textWidth = textHeight*alpha
	var imageWidth = int(float64(textXCount)*textWidth)
	img:=image.NewRGBA(image.Rect(0,0,imageWidth,imageHeight))
	for j:=0;j<imageHeight;j++{
		for i:=0;i<imageWidth;i++{
			img.Set(i,j,color.RGBA{0,0,0,255})
		}
	}
	font,err:=freetype.ParseFont(fontBytes)
	if err!=nil{
		return nil,err
	}
	c:=freetype.NewContext()
	c.SetClip(img.Bounds())
	c.SetFont(font)
	c.SetSrc(image.White)
	c.SetDst(img)
	c.SetFontSize(textHeight)
	for y:=0;y<textYCount;y++{
		for x:=0;x<textXCount;x++{
			c.DrawString(string(text[y*textXCount+x]),freetype.Pt(int(float64(x)*textWidth),int(float64(y)*textHeight+textHeight)))
		}
	}
	return img,nil
}