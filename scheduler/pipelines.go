package scheduler

import (
	"fmt"
	"proj1/png"
	"strings"
	"sync"
)

type ProcessingTask struct {
	InputPath  string
	OutputPath string
	Effects    []string
	Dir        string
}

type TaskedImage struct {
	Image *png.Image
	Task  ProcessingTask
}

func RunPipeline(config Config) {
	if config.ThreadCount == -1 {
		fmt.Println("Performance test: Execution ended")
		return
	}

	effectsData, err := getEffectData()
	if err != nil {
		fmt.Println("Failed to get effects data from effects.txt:", err)
		return
	}

	dirs := strings.Split(config.DataDirs, "+")

	taskChan := make(chan ProcessingTask, config.ThreadCount)
	var wg sync.WaitGroup

	chunkSize := (len(dirs) + config.ThreadCount - 1) / config.ThreadCount
	for i := 0; i < len(dirs); i += chunkSize {
		end := i + chunkSize
		if end > len(dirs) {
			end = len(dirs)
		}
		wg.Add(1)
		go func(subDirs []string) {
			defer wg.Done()
			for _, dir := range subDirs {
				for _, effect := range effectsData {
					taskChan <- ProcessingTask{
						Dir:        dir,
						InputPath:  effect.InputPath,
						OutputPath: effect.OutputPath,
						Effects:    effect.Effects,
					}
				}
			}
		}(dirs[i:end])
	}

	go func() {
		wg.Wait()
		close(taskChan)
	}()

	TaskChan := make([][]ProcessingTask, config.ThreadCount)
	index := 0
	for task := range taskChan {
		TaskChan[index] = append(TaskChan[index], task)
		index = (index + 1) % config.ThreadCount
	}

	var wgProcess sync.WaitGroup
	imgChan := make(chan TaskedImage, config.ThreadCount)
	for i := 0; i < config.ThreadCount; i++ {
		wgProcess.Add(1)
		go func(batch []ProcessingTask) {
			defer wgProcess.Done()
			for _, task := range batch {
				img := loadImage(task.Dir, task.InputPath)
				if img != nil {
					img.ProcessImage(task.Dir, task.OutputPath, task.Effects, img.Bounds.Min.Y, img.Bounds.Max.Y, true)
					imgChan <- TaskedImage{Image: img, Task: task}
				}
			}
		}(TaskChan[i])
	}

	go func() {
		wgProcess.Wait()
		close(imgChan)
	}()

	var wgSave sync.WaitGroup
	for taskedImg := range imgChan {
		wgSave.Add(1)
		go func(img TaskedImage) {
			defer wgSave.Done()
			saveImage(img.Image, img.Task.Dir, img.Task.OutputPath)
		}(taskedImg)
	}
	wgSave.Wait()
}
