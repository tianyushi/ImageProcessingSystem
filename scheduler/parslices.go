package scheduler

import (
	"fmt"
	"strings"
	"sync"
)

func RunParallelSlices(config Config) {

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
				sliceLength := img.Bounds.Dy() / threadNum
				for _, effect := range task.Effects {
					var sliceWg sync.WaitGroup
					for j := 0; j < threadNum; j++ {
						sliceWg.Add(1)
						start := j * sliceLength
						end := (j + 1) * sliceLength
						if j == threadNum-1 {
							end = img.Bounds.Max.Y
						}

						overlap := 2
						if start > 0 {
							start -= overlap
						}
						if end < img.Bounds.Max.Y {
							end += overlap
						}

						go func(s, e int) {
							defer sliceWg.Done()
							img.ProcessImage(task.Dir, task.OutputPath, []string{effect}, s, e, false)
						}(start, end)
					}
					sliceWg.Wait()
					img.SwapInAndOut()
				}
				img.SwapInAndOut()
				saveImage(img, task.Dir, task.OutputPath)
			}
		}()
	}

	wg.Wait()
}
