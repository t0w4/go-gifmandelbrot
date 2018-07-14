// マンデルブロフラクタルのGIF画像を生成する.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/gif"
	"math/cmplx"
	"os"
	"sync"
)

var palette = []color.Color{color.White, color.Black}

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
	nframes                = 200 // アニメーションフレーム数
	delay                  = 8   // 10ms単位でのフレーム間の遅延
)

var powArg float64

func init()  {
	flag.Float64Var(&powArg,"pow", 2.0, "ext) ./go-gifmandelbrot -pow=2.0")
}

func main() {
	flag.Parse()
	anim := gif.GIF{LoopCount: nframes}
	for pow := 1.0; pow < powArg; pow += 0.02 {
		draw(&anim, pow)
	}
	// gif に保存する
	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &anim)
}

func mandelbrot(z complex128, pow float64) color.Color {
	const iterations = 200
	const contrast = 15
	var v complex128 = 0
	for n := uint8(0); n < iterations; n++ {
		v = cmplx.Pow(v, complex(pow, 0)) + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{Y: 255 - contrast*n}
		}
	}
	return color.Black
}

func draw(anim *gif.GIF, pow float64) {
	rect := image.Rect(0, 0, width, height)
	img := image.NewPaletted(rect, palette)

	var wg sync.WaitGroup
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			x := px
			y := py
			wg.Add(1)
			go func() {
				defer wg.Done()
				plot(img, x, y, pow)
			}()
		}
	}
	wg.Wait()
	anim.Delay = append(anim.Delay, delay)
	anim.Image = append(anim.Image, img)
}

func plot(img *image.Paletted, px int, py int, pow float64)  {
	y := float64(py)/height*(ymax-ymin) + ymin
	x := float64(px)/width*(xmax-xmin) + xmin
	z := complex(x, y)
	// 画像の点(px, py)は複素数値zを表している。
	img.Set(px, py, mandelbrot(z, pow))
}