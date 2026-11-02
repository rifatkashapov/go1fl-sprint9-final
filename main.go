package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

func applyFuncToChunk(wg *sync.WaitGroup, size int, fn func(*sync.WaitGroup, int, int, int)) {
	if wg == nil || size <= 0 || fn == nil {
		return
	}
	chunks := CHUNKS
	if size < CHUNKS {
		chunks = size
	}

	itemsInChunk := size / chunks
	remainder := size % chunks

	chunkStartPos := 0
	for chunkNumber := 0; chunkNumber < chunks; chunkNumber++ {
		wg.Add(1)
		chunkSize := itemsInChunk
		if remainder > 0 {
			chunkSize++
			remainder--
		}
		go fn(wg, chunkNumber, chunkStartPos, chunkStartPos+chunkSize)
		chunkStartPos += chunkSize
	}

	wg.Wait()
}

// generateRandomElements returns positive random integers, or an empty slice for size <= 0.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}

	var result = make([]int, size)
	var wg sync.WaitGroup

	generate := func(wg *sync.WaitGroup, _ int, start, stop int) {
		defer wg.Done()
		for i := start; i < stop; i++ {
			result[i] = rand.Intn(SIZE) + 1
		}
	}

	applyFuncToChunk(&wg, size, generate)

	return result
}

// maximum returns the largest element, or 0 for an empty slice.
func maximum(data []int) int {
	return findMax(data, 0, len(data))
}

// findMax returns the maximum in [start, stop), or 0 for an invalid range.
func findMax(data []int, start, stop int) int {
	if start < 0 || start >= stop || stop > len(data) {
		return 0
	}

	initialMax := data[start]
	for i := start + 1; i < stop; i++ {
		if data[i] > initialMax {
			initialMax = data[i]
		}
	}
	return initialMax
}

// maxChunks finds the largest element concurrently, or returns 0 for an empty slice.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}

	result := make([]int, min(CHUNKS, len(data)))

	var wgExt sync.WaitGroup

	findMaxNum := func(wgInt *sync.WaitGroup, chunkNum int, start, stop int) {
		defer wgInt.Done()
		result[chunkNum] = findMax(data, start, stop)
	}

	applyFuncToChunk(&wgExt, len(data), findMaxNum)
	return findMax(result, 0, len(result))
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	arr := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now()
	maxNumber := maximum(arr)
	elapsed := time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %v us\n", maxNumber, elapsed)

	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)
	start = time.Now()
	maxNumber = maxChunks(arr)
	elapsed = time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %v us\n", maxNumber, elapsed)
}
