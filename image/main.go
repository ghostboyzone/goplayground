package main

import (
	// "fmt"
	"image"
	// "image/jpeg"
	"github.com/disintegration/imaging"
	"log"
	// "os"
	"image/color"
)

type ImageHash struct {
	Rate     float64
	AvgColor color.Color
}

func main() {
	imgIn, err := imaging.Open("test_4.png")
	if err != nil {
		log.Fatal(err)
	}
	var imgHash ImageHash
	imgRec := imgIn.Bounds()
	log.Println("imgRec", imgRec)
	imgHash.Rate = (float64)(imgRec.Max.X) / (float64)(imgRec.Max.Y)
	imgDst := imaging.New(imgRec.Max.X, imgRec.Max.Y, color.NRGBA{0, 0, 0, 0})
	imgDst = imaging.Paste(imgDst, imgIn, image.Pt(0, 0))

	imaging.Save(imgDst, "tmp.png")
	log.Println(imgHash)

	totalPixes := (float64)(imgRec.Max.X) * (float64)(imgRec.Max.Y)
	log.Println(totalPixes)
	var avgR, avgG, avgB, avgA float64
	avgR = 0
	for i := 0; i < imgRec.Max.X; i++ {
		for j := 0; j < imgRec.Max.Y; j++ {
			r, g, b, a := imgDst.At(i, j).RGBA()
			avgR += (float64)(r) / totalPixes
			if (float64)(r)/totalPixes > 0 {

				log.Fatal(r, g, b, a, avgR)
			}
			avgG += (float64)(g) / totalPixes
			avgB += (float64)(b) / totalPixes
			avgA += (float64)(a) / totalPixes
		}
	}
	log.Fatal(avgR, avgG, avgB, avgA)
}
