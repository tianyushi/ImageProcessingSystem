# High Performance Image Processing System Implemented in Go

## Overview
This project implements an image editor that applies image effects using 2D convolutions on a series of images. The project involves three implementations:
1. A sequential baseline version.
2. A parallel version that processes multiple images concurrently.
3. A parallel version that processes slices of each image concurrently.

## Image Effects
The following effects are supported:
- **Sharpen (S)**: Uses a convolution kernel `[0, -1, 0, -1, 5, -1, 0, -1, 0]`.
- **Edge Detection (E)**: Uses a convolution kernel `[-1, -1, -1, -1, 8, -1, -1, -1, -1]`.
- **Blur (B)**: Uses a convolution kernel `[1/9, 1/9, 1/9, 1/9, 1/9, 1/9, 1/9, 1/9, 1/9]`.
- **Grayscale (G)**: Averages the RGB values of each pixel.

## JSON Input Format
The program reads JSON strings with the following format:
```json
{
  "inPath": "sky.png",
  "outPath": "sky_out.png",
  "effects": ["S", "B", "E"]
}
```
- **inPath**: Input image file path.
- **outPath**: Output image file path.
- **effects**: List of effects to apply in order.

## Directory Structure
- **effects.txt**: Contains the JSON strings.
- **expected**: Contains the expected output images.
- **in**: Contains subdirectories `big`, `mixture`, and `small` with input images.
- **out**: Directory for saving output images.

## Implementations

### 1. Sequential Implementation
Implemented in `proj1/scheduler/sequential.go`:
```go
func RunSequential(config Config) {
    // Implementation code
}
```

### 2. Multiple Images in Parallel
Implemented in `proj1/scheduler/parallel_images.go`:
- Creates a task queue for images.
- Uses Go routines to process images concurrently.
- Employs a TAS lock for thread safety.

### 3. Parallelize Each Image
Implemented in `proj1/scheduler/parallel_slices.go`:
- Spawns Go routines to process image slices concurrently.
- Manages dependencies between effects using waitgroups or overlapping slices.

### 4. Pipeline Parallelization
Implemented in `proj1/scheduler/pipeline.go`:
- Divides tasks into pipelines to process different stages of image processing concurrently.
- Uses channels to communicate between different stages of the pipeline.
- Balances workload dynamically across multiple Go routines.

### 5. Work Stealing Parallelization
Implemented in `proj1/scheduler/work_stealing.go`:
- Utilizes a work-stealing algorithm to balance the load dynamically across multiple workers.
- Implements concurrent deques for each worker to manage tasks.
- Allows workers to steal tasks from others when they run out of their own tasks.

## Performance Measurements
- Measure performance using Linux cluster tools.
- Calculate speedup using Amdahl's law for various thread counts.
- Create speedup graphs for the sequential and parallel implementations.

## Performance Analysis Report
Includes:
- Project description.
- Instructions for running the testing script.
- Analysis of performance graphs and results.
- Answers to provided questions on hotspots, bottlenecks, and performance evaluation.

## Running the Program
To run the program, use the following command:
```sh
go run editor.go <data_dir> <mode> <number_of_threads>
```
- `<data_dir>`: Specifies the data directory (e.g., `big`, `small`, `mixture`).
- `<mode>`: Specifies the mode (e.g., `sequential`, `parallel_images`, `parallel_slices`).
- `<number_of_threads>`: Specifies the number of threads for parallel modes.

## Example Commands
Sequential:
```sh
go run editor.go big
```
Parallel Images:
```sh
go run editor.go big parallel_images 4
```
Parallel Slices:
```sh
go run editor.go big parallel_slices 4
```

## Conclusion
This project demonstrates the use of parallel programming to enhance the performance of an image processing system, applying effects using 2D convolutions. The implementations showcase the benefits and challenges of parallelization in a CPU-based environment.
