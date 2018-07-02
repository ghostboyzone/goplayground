package main

import (
	"fmt"
	"image"
	// "image/jpeg"
	"github.com/disintegration/imaging"
	"log"
	// "os"
	"container/list"
	"image/color"
	"math"
	// "sort"
	"strconv"
)

type ImageHash struct {
	Ratio    float64
	AvgColor color.NRGBA
	MapData  XYMap
}

func (imgHash ImageHash) String() string {
	return fmt.Sprintf("Ratio[%f] AvgColor[%+v] MapData_Count[%d]", imgHash.Ratio, imgHash.AvgColor, len(imgHash.MapData))
}

func (a ImageHash) CalDiff(b ImageHash) float64 {
	if len(a.MapData) != len(b.MapData) {
		return -1.0
	}
	maxIdx := len(a.MapData)
	var sum1, sum2, sum3 float64
	sum1 = 0
	sum2 = 0
	sum3 = 0
	for i := 0; i < maxIdx; i++ {
		sum1 += a.MapData[i] * b.MapData[i]
		sum2 += a.MapData[i] * a.MapData[i]
		sum3 += b.MapData[i] * b.MapData[i]
	}
	return sum1 / (math.Sqrt(sum2) * math.Sqrt(sum3))
}

type SmallRGB struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type SmallRGBCounter struct {
	Color SmallRGB
	Count int64
}

type SmallRGBMap map[string]SmallRGBCounter

//x-y-r-g-b-a
//x,y (0,7)
//r,g,b,a (0,7)
type XYMap map[int]float64

/*
type IntSlice []int

func (c IntSlice) Len() int {
	return len(c)
}
func (c IntSlice) Less(i, j int) int {
	return c[i] < c[j]
}

func (c IntSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
*/

var (
	smallRGBMap SmallRGBMap
)

func main() {
	InitRGBMap()
	imgIn, err := imaging.Open("001.png")
	if err != nil {
		log.Fatal(err)
	}
	imgHash := GenImageHash(imgIn)
	log.Println(imgHash)

	imgIn1, err := imaging.Open("002.png")
	if err != nil {
		log.Fatal(err)
	}
	imgHash1 := GenImageHash(imgIn1)
	log.Println(imgHash1)
	log.Println(imgHash.CalDiff(imgHash1))
}

func InitRGBMap() {
	smallRGBMap = make(SmallRGBMap)
	resultList := getAllRGBMap("", 0, 4)
	fmt.Println(resultList.Len())
	for e := resultList.Front(); e != nil; e = e.Next() {
		smallRGBMap[e.Value.(string)] = SmallRGBCounter{
			Color: parseSmallRGB(e.Value.(string)),
			Count: 0,
		}
	}
}

func GenImageHash(img image.Image) ImageHash {
	var imgHash ImageHash
	imgRec := img.Bounds()
	log.Println("imgRec", imgRec)
	imgHash.Ratio = (float64)(imgRec.Max.X) / (float64)(imgRec.Max.Y)
	imgDst := imaging.New(imgRec.Max.X, imgRec.Max.Y, color.NRGBA{0, 0, 0, 0})
	imgDst = imaging.Paste(imgDst, img, image.Pt(0, 0))

	imaging.Save(imgDst, "tmp.png")

	totalPixes := (float64)(imgRec.Max.X) * (float64)(imgRec.Max.Y)
	log.Println(totalPixes)
	var avgR, avgG, avgB, avgA float64
	for i := 0; i < imgRec.Max.X; i++ {
		for j := 0; j < imgRec.Max.Y; j++ {
			nrgbaTmp := imgDst.NRGBAAt(i, j)
			avgR += (float64)(nrgbaTmp.R) / totalPixes
			avgG += (float64)(nrgbaTmp.G) / totalPixes
			avgB += (float64)(nrgbaTmp.B) / totalPixes
			avgA += (float64)(nrgbaTmp.A) / totalPixes
		}
	}

	imgHash.AvgColor = color.NRGBA{
		R: uint8(math.Floor(avgR)),
		G: uint8(math.Floor(avgG)),
		B: uint8(math.Floor(avgB)),
		A: uint8(math.Floor(avgA)),
	}

	xyMap := make(XYMap)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			for k, _ := range smallRGBMap {
				// x-y-r-g-b-a
				// 15  12  9  6 3  0

				smallTmp := parseSmallRGB(k)

				tmpKey := i<<15 + j<<12 + (int)(smallTmp.R)<<9 + (int)(smallTmp.G)<<6 + (int)(smallTmp.B)<<3 + (int)(smallTmp.A)

				xyMap[tmpKey] = 0
			}
		}
	}

	for i := 0; i < imgRec.Max.X; i++ {
		for j := 0; j < imgRec.Max.Y; j++ {
			nrgbaTmp := imgDst.NRGBAAt(i, j)

			myX := (int)(math.Ceil((float64)(i) / (math.Ceil((float64)(imgRec.Max.X) / 8.0))))
			if myX > 0 {
				myX -= 1
			}
			myY := (int)(math.Ceil((float64)(j) / (math.Ceil((float64)(imgRec.Max.Y) / 8.0))))
			if myY > 0 {
				myY -= 1
			}

			smallTmp := transNRGBA2SmallRGB(nrgbaTmp)

			tmpKey := myX<<15 + myY<<12 + (int)(smallTmp.R)<<9 + (int)(smallTmp.G)<<6 + (int)(smallTmp.B)<<3 + (int)(smallTmp.A)

			xyMap[tmpKey]++
		}
	}

	imgHash.MapData = xyMap
	// for idx := 0; idx < 262144; idx++ {
	// 	// log.Fatal(idx, xyMap[idx])
	// }
	// log.Println(imgHash)
	return imgHash
}

func getAllRGBMap(v string, depth int, maxDepth int) *list.List {
	if depth == maxDepth {
		l := list.New()
		l.PushFront(v)
		return l
	}
	l := list.New()
	for i := 0; i < 8; i++ {
		result := getAllRGBMap(v+strconv.Itoa(i), depth+1, maxDepth)
		for e := result.Front(); e != nil; e = e.Next() {
			l.PushFront(e.Value)
		}
	}
	return l
}

func parseSmallRGB(c string) SmallRGB {
	s := SmallRGB{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	}
	if len(c) != 4 {
		return s
	}
	var tmp int
	tmp, _ = strconv.Atoi(c[0:1])
	s.R = (uint8)(tmp)
	tmp, _ = strconv.Atoi(c[1:2])
	s.G = (uint8)(tmp)
	tmp, _ = strconv.Atoi(c[2:3])
	s.B = (uint8)(tmp)
	tmp, _ = strconv.Atoi(c[3:4])
	s.A = (uint8)(tmp)
	return s
}

func transNRGBA2SmallRGB(c color.NRGBA) SmallRGB {
	return SmallRGB{
		R: c.R >> 5,
		G: c.G >> 5,
		B: c.B >> 5,
		A: c.A >> 5,
	}
}
