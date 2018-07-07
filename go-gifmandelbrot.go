// マンデルブロフラクタルのGIF画像を生成する.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"math/cmplx"
	"os"
)

var palette = []color.Color{color.White, color.Black}

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
	nframes                = 200 // アニメーションフレーム数
	delay                  = 8   // 10ms単位でのフレーム間の遅延
)

func main() {
	anim := gif.GIF{LoopCount: nframes}
	for pow := 1.0; pow < 10.0; pow += 0.02 {
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

	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// 画像の点(px, py)は複素数値zを表している。
			img.Set(px, py, mandelbrot(z, pow))
		}
	}
	anim.Delay = append(anim.Delay, delay)
	anim.Image = append(anim.Image, img)
}