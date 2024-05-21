package scheduler

import (
	"fmt"
	"strings"
)

func RunSequential(config Config) {
	directories := strings.Split(config.DataDirs, "+")
	effectData, err := getEffectData()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, effect := range effectData {
		for _, dir := range directories {
			img := loadImage(dir, effect.InputPath)
			img.ProcessImage(dir, effect.OutputPath, effect.Effects, img.Bounds.Min.Y, img.Bounds.Max.Y, true)
			saveImage(img, dir, effect.OutputPath)
		}
	}

}
