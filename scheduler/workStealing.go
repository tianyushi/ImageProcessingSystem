package scheduler

import (
	"math/rand"
	"strings"
	"sync"
)

type ImageProcessingTask struct {
	Dir        string
	InputPath  string
	OutputPath string
	Effects    []string
}

func worker(ConcurrentDeques []*ConcurrentDeque, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {

		task := ConcurrentDeques[id].PopBottom()
		if task == nil {
			task = stealTask(ConcurrentDeques, id)
			if task == nil {
				return
			}
		}

		processTask(task.(ImageProcessingTask))
	}
}

func stealTask(ConcurrentDeques []*ConcurrentDeque, stealed int) interface{} {
	n := len(ConcurrentDeques)
	start := rand.Intn(n)
	for i := 0; i < n; i++ {
		victim := (start + i) % n
		if victim == stealed {
			continue
		}
		if task := ConcurrentDeques[victim].Steal(); task != nil {
			return task
		}
	}
	return nil
}

func processTask(task ImageProcessingTask) {
	img := loadImage(task.Dir, task.InputPath)
	if img != nil {
		img.ProcessImage(task.Dir, task.OutputPath, task.Effects, img.Bounds.Min.Y, img.Bounds.Max.Y, true)
		saveImage(img, task.Dir, task.OutputPath)
	}
}

func RunWorkStealing(config Config) {
	dirs := strings.Split(config.DataDirs, "+")
	effectData, err := getEffectData()

	if err != nil {
		panic(err)
	}

	ConcurrentDeques := make([]*ConcurrentDeque, config.ThreadCount)
	for i := range ConcurrentDeques {
		ConcurrentDeques[i] = NewDequeue(8)
	}

	for _, data := range effectData {
		for _, dir := range dirs {
			task := ImageProcessingTask{
				Dir:        dir,
				InputPath:  data.InputPath,
				OutputPath: data.OutputPath,
				Effects:    data.Effects,
			}
			ConcurrentDeques[rand.Intn(len(ConcurrentDeques))].PushBottom(task)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < config.ThreadCount; i++ {
		wg.Add(1)
		go worker(ConcurrentDeques, i, &wg)
	}
	wg.Wait()
}
