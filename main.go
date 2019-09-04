package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/draw"

	"github.com/urfave/cli"
)

func main() {
	// 引数を取得
	srcImagePath := flag.String("src_path", "", "source image path")
	dstImagePath := flag.String("dst_path", "", "destination image path")
	dstFormat := flag.String("dst_format", "", "destination image format. e.g) png,jpeg,gif")
	height := flag.Int("height", 0, "height size")
	width := flag.Int("width", 0, "width size")
	flag.Parse()

	if *srcImagePath == "" || *dstImagePath == "" || *dstFormat == "" || *height == 0 || *width == 0 {
		help()
		return
	}

	// srcファイルを開く
	srcImgFile, err := os.Open(*srcImagePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer srcImgFile.Close()

	// 画像読み込み
	srcImg, _, err := image.Decode(srcImgFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	srcRct := srcImg.Bounds()
	dstImg := image.NewRGBA(image.Rect(0, 0, *width, *height))
	dstRct := dstImg.Bounds()
	draw.CatmullRom.Scale(dstImg, dstRct, srcImg, srcRct, draw.Over, nil)
	dstImgFile, err := os.Create(*dstImagePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	defer dstImgFile.Close()

	// 指定の拡張子にエンコード
	switch *dstFormat {
	case "jpeg":
		err := jpeg.Encode(dstImgFile, dstImg, &jpeg.Options{Quality: 100})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "gif":
		err := gif.Encode(dstImgFile, dstImg, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "png":
		err := png.Encode(dstImgFile, dstImg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func help() {
	app := cli.NewApp()
	app.Name = "image_resizer"
	app.Usage = "this is image size converter. "
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "src_path",
			Usage: "Load source image file path from `FILE`",
		},
		cli.StringFlag{
			Name:  "dst_path",
			Usage: "Load destination image file path from `FILE`",
		},
		cli.StringFlag{
			Name:  "dst_format",
			Usage: "Load destination image file extension value",
		},
		cli.StringFlag{
			Name:  "height",
			Usage: "height size",
		},
		cli.StringFlag{
			Name:  "width",
			Usage: "width size",
		},
	}
	app.Run(os.Args)
	return
}
