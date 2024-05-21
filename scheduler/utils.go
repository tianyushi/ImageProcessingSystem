package scheduler

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"proj1/png"
)

type EffectData struct {
	InputPath  string   `json:"inPath"`
	OutputPath string   `json:"outPath"`
	Effects    []string `json:"effects"`
}

func getEffectData() ([]EffectData, error) {
	filePath := "../data/effects.txt"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening effects.txt file %s: %w", filePath, err)
	}
	defer file.Close()

	var effects []EffectData
	decoder := json.NewDecoder(file)
	for {
		var effect EffectData
		err := decoder.Decode(&effect)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error decoding JSON: %w", err)
		}
		effects = append(effects, effect)
	}
	return effects, nil
}

func loadImage(dir, InputPath string) *png.Image {
	inputPath := fmt.Sprintf("../data/in/%s/%s", dir, InputPath)
	img, err := png.Load(inputPath)
	if err != nil {
		fmt.Println("Error loading the image:", err)
		return nil
	}

	return img
}

func saveImage(img *png.Image, dir, filename string) {
	outputPath := fmt.Sprintf("../data/out/%s_%s", dir, filename)
	if err := img.Save(outputPath); err != nil {
		fmt.Println("Error saving the image:", err)
	}
}

func createTasks(dirs []string, effectsData []EffectData) []ProcessingTask {
	var tasks []ProcessingTask
	for _, dir := range dirs {
		for _, effect := range effectsData {
			tasks = append(tasks, ProcessingTask{
				Dir:        dir,
				InputPath:  effect.InputPath,
				OutputPath: effect.OutputPath,
				Effects:    effect.Effects,
			})
		}
	}
	return tasks
}
