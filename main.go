package main

import (
	"flag"
	"math"

	"github.com/fogleman/gg"
	"github.com/garciaolais/mlx90640"
)

func main() {
	tempMinPtr := flag.Float64("min", 30.0, "temp min")
	tempMaxPtr := flag.Float64("max", 38.0, "temp max")
	flag.Parse()

	m := mlx90640.New("/dev/i2c-1")
	defer m.Close()
	m.SetSampling(mlx90640.Sampling_02_Hz)

	dc := gg.NewContext(mlx90640.WIDTH, mlx90640.HEIGHT)
	frame := m.GetFrame()
	for h := 0; h < mlx90640.HEIGHT; h++ {
		for w := 0; w < mlx90640.WIDTH; w++ {
			t := frame[h*mlx90640.WIDTH+w]
			tempNormalized := (t - *tempMinPtr) / (*tempMaxPtr - *tempMinPtr)
			red, green, blue := GetRGBColor(tempNormalized)
			dc.DrawPoint(float64(w), float64(h), 1)
			dc.SetRGB255(int(red), int(green), int(blue))
			dc.Fill()
		}
	}

	dc.SavePNG("out.png")
}

func GetRGBColor(v float64) (red, green, blue int) {
	var idx1, idx2 int
	var fractBetween float64

	colorSpectrum := [][]float64{
		{0, 0, 127},
		{0, 0, 255},
		{0, 255, 0},
		{255, 255, 0},
		{255, 0, 0},
	}

	if v <= 0.0 {
		idx1 = 0
		idx2 = 0
	} else if v >= 1.0 {
		idx1 = len(colorSpectrum) - 1
		idx2 = len(colorSpectrum) - 1
	} else {
		v = v * (float64(len(colorSpectrum)) - 1)
		idx1 = int(math.Floor(float64(v)))
		idx2 = idx1 + 1
		fractBetween = v - float64(idx1)
	}

	var ir, ig, ib int

	ir = (int)((((colorSpectrum[idx2][0] - colorSpectrum[idx1][0]) * fractBetween) + colorSpectrum[idx1][0]))
	ig = (int)((((colorSpectrum[idx2][1] - colorSpectrum[idx1][1]) * fractBetween) + colorSpectrum[idx1][1]))
	ib = (int)((((colorSpectrum[idx2][2] - colorSpectrum[idx1][2]) * fractBetween) + colorSpectrum[idx1][2]))

	return ir, ig, ib
}
