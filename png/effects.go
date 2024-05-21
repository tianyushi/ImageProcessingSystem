// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
)

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale(start, end int) {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()

	if start < bounds.Min.Y {
		start = bounds.Min.Y
	}
	if end > bounds.Max.Y {
		end = bounds.Max.Y
	}
	for y := start; y < end; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

func (img *Image) Convolution(kernel [][]float64, start, end int) {
	kernelCenter := len(kernel) / 2
	bounds := img.out.Bounds()

	originalStart := start
	originalEnd := end

	start = start - kernelCenter
	end = end + kernelCenter

	if start < bounds.Min.Y {
		start = bounds.Min.Y
	}
	if end > bounds.Max.Y {
		end = bounds.Max.Y
	}

	for y := start; y < end; y++ {
		if y < originalStart || y >= originalEnd {
			continue
		}
		for x := bounds.Min.X + kernelCenter; x < bounds.Max.X-kernelCenter; x++ {
			var sumR, sumG, sumB float64
			_, _, _, a := img.in.At(x, y).RGBA()
			for ky := 0; ky < len(kernel); ky++ {
				for kx := 0; kx < len(kernel); kx++ {
					nx, ny := x+kx-kernelCenter, y+ky-kernelCenter
					r, g, b, _ := img.in.At(nx, ny).RGBA()
					kernelVal := kernel[ky][kx]
					sumR += kernelVal * float64(r)
					sumG += kernelVal * float64(g)
					sumB += kernelVal * float64(b)
				}
			}
			img.out.Set(x, y, color.RGBA64{
				R: clamp(sumR),
				G: clamp(sumG),
				B: clamp(sumB),
				A: uint16(a),
			})
		}
	}
}

func (img *Image) Sharpen(start, end int) {
	sharpenKernel := [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
	img.Convolution(sharpenKernel, start, end)
}

func (img *Image) EdgeDetection(start, end int) {
	edgeKernel := [][]float64{
		{-1, -1, -1},
		{-1, 8, -1},
		{-1, -1, -1},
	}
	img.Convolution(edgeKernel, start, end)
}

func (img *Image) Blur(start, end int) {
	blurKernel := [][]float64{
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
		{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
	}
	img.Convolution(blurKernel, start, end)
}

func (img *Image) SwapInAndOut() {
	img.in, img.out = img.out, img.in
}

func (img *Image) ProcessImage(dir, outputPath string, effects []string, start, end int, swapFlag bool) {
	for _, effect := range effects {
		switch effect {
		case "G":
			img.Grayscale(start, end)
		case "B":
			img.Blur(start, end)
		case "S":
			img.Sharpen(start, end)
		case "E":
			img.EdgeDetection(start, end)
		}

		if swapFlag {
			img.SwapInAndOut()
		}
	}

	if swapFlag {
		img.SwapInAndOut()
	}
}
