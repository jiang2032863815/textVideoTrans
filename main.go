package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"
	"videoToTextVideo/imageToText"
	"videoToTextVideo/textImage"
)

var inputVideoPath = flag.String("i","","inputVideoPath")
var outputVideoHeight = flag.Int("h",0,"outputVideoHeight")
var textXCount = flag.Int("x",0,"textXCount")
var textYCount = flag.Int("y",0,"textYCount")
//go:embed f.ttf
var fontBytes []byte
func checkFlags(){
	if *inputVideoPath==""{
		log.Fatalln("-i 请指定mp4文件路径")
	}
	if *outputVideoHeight==0{
		log.Fatalln("-h 请指定输出视频高度")
	}
	if *textXCount==0{
		log.Fatalln("-x 请指定横坐标方向用多少个字符")
	}
	if *textYCount==0{
		log.Fatalln("-y 请指定纵坐标方向用多少个字符")
	}
}
func mkdir(name string){
	oldMask:=syscall.Umask(0)
	os.Mkdir(name,os.ModePerm)
	syscall.Umask(oldMask)
}
func main(){
	flag.Parse()
	checkFlags()
	var dirName=strconv.FormatInt(time.Now().UnixNano(),10)+"_tmp"
	mkdir(dirName)
	cmd:=exec.Command("ffmpeg","-i",*inputVideoPath,"-r","60","-ss","00:00:00",dirName+"/%d.png")
	err:=cmd.Run()
	if err!=nil{
		log.Fatalln(err)
	}
	info,err:=ioutil.ReadDir(dirName)
	if err!=nil{
		log.Fatalln(err)
	}
	outDirName:=strconv.FormatInt(time.Now().UnixNano(),10)+"__tmp"
	mkdir(outDirName)
	for _,v:=range info{
		if path.Ext(v.Name())!=".png"{
			continue
		}
		text,err:=imageToText.ImgToTxt(path.Join(dirName,v.Name()),*textXCount,*textYCount,"",imageToText.DefaultTxt)
		if err!=nil{
			log.Fatalln(err)
		}
		rgba,err:=textImage.DrawTextToImageRect(text,*textXCount,*textYCount,*outputVideoHeight,fontBytes)
		if err!=nil{
			log.Fatalln(err)
		}
		outputFile,err:=os.Create(path.Join(outDirName,v.Name()))
		if err!=nil{
			log.Fatalln(err)
		}
		png.Encode(outputFile,rgba)
		outputFile.Close()
		fmt.Println(v.Name())
	}
	cmd=exec.Command("ffmpeg","-r","60","-i",outDirName+"/%d.png","out_"+*inputVideoPath)
	err=cmd.Run()
	if err!=nil{
		log.Fatalln(err)
	}
	os.RemoveAll(dirName)
	os.RemoveAll(outDirName)
}