package scheduler

import (
	"fmt"
	"strings"
	"sync"
)

func RunParallelFiles(config Config) {

	if config.ThreadCount == -1 {
		fmt.Println("performance test: excution eneded")
		return
	}

	effectsData, err := getEffectData()
	if err != nil {
		fmt.Println("Failed to get effects data from effects.txt:", err)
		return
	}

	dirs := strings.Split(config.DataDirs, "+")
	tasks := createTasks(dirs, effectsData)

	var wg sync.WaitGroup
	lock := NewTASLock{}

	threadNum := config.ThreadCount
	if len(tasks) < threadNum {
		threadNum = len(tasks)
	}

	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				lock.Lock()
				if len(tasks) == 0 {
					lock.Unlock()
					return
				}
				task := tasks[0]
				tasks = tasks[1:]
				lock.Unlock()
				img := loadImage(task.Dir, task.InputPath)
				img.ProcessImage(task.Dir, task.OutputPath, task.Effects, img.Bounds.Min.Y, img.Bounds.Max.Y, true)
				saveImage(img, task.Dir, task.OutputPath)
			}
		}()
	}

	wg.Wait()
}
