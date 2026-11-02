package main

import (
	"math"
	"slices"
	"sync"
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantLen int
	}{
		{name: "отрицательный размер", size: -1, wantLen: 0},
		{name: "минимальный int", size: math.MinInt, wantLen: 0},
		{name: "нулевой размер", size: 0, wantLen: 0},
		{name: "один элемент", size: 1, wantLen: 1},
		{name: "меньше числа частей", size: CHUNKS - 1, wantLen: CHUNKS - 1},
		{name: "ровно восемь частей", size: CHUNKS, wantLen: CHUNKS},
		{name: "размер с остатком", size: CHUNKS + 1, wantLen: CHUNKS + 1},
		{name: "большой слайс", size: 100_003, wantLen: 100_003},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRandomElements(tt.size)
			if len(got) != tt.wantLen {
				t.Fatalf("generateRandomElements(%d): длина = %d, ожидается %d", tt.size, len(got), tt.wantLen)
			}
			for i, value := range got {
				if value <= 0 {
					t.Fatalf("generateRandomElements(%d): элемент [%d] = %d, ожидается положительное число", tt.size, i, value)
				}
			}
		})
	}
}

// Общая таблица проверяет обе реализации по одним и тем же ожидаемым значениям.
func TestMaximum(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want int
	}{
		{name: "nil", data: nil, want: 0},
		{name: "пустой слайс", data: []int{}, want: 0},
		{name: "один положительный", data: []int{5}, want: 5},
		{name: "один отрицательный", data: []int{-5}, want: -5},
		{name: "один ноль", data: []int{0}, want: 0},
		{name: "максимум в начале", data: []int{9, 3, 2, 1}, want: 9},
		{name: "максимум в середине", data: []int{1, 2, 9, 3, 4}, want: 9},
		{name: "максимум в конце", data: []int{1, 2, 3, 9}, want: 9},
		{name: "одинаковые значения", data: []int{7, 7, 7}, want: 7},
		{name: "нули", data: []int{0, 0, 0}, want: 0},
		{name: "повторяющийся максимум", data: []int{9, 1, 9, 2, 9}, want: 9},
		{name: "короткий отрицательный слайс", data: []int{-9, -2, -7}, want: -2},
		{name: "смешанные знаки", data: []int{-9, 0, 4, -1}, want: 4},
		{name: "ровно восемь элементов", data: []int{1, 2, 3, 4, 5, 6, 7, 8}, want: 8},
		{name: "максимум в остатке", data: []int{1, 2, 3, 4, 5, 6, 7, 8, 99}, want: 99},
		{name: "несколько элементов в части", data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, want: 16},
		{name: "отрицательные с остатком", data: []int{-9, -8, -7, -6, -5, -4, -3, -2, -1}, want: -1},
		{name: "минимальный int", data: []int{math.MinInt}, want: math.MinInt},
		{name: "границы int", data: []int{math.MinInt, 0, math.MaxInt}, want: math.MaxInt},
	}
	functions := []struct {
		name string
		fn   func([]int) int
	}{
		{name: "maximum", fn: maximum},
		{name: "maxChunks", fn: maxChunks},
	}
	for _, function := range functions {
		t.Run(function.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					data := slices.Clone(tt.data)
					if got := function.fn(data); got != tt.want {
						t.Errorf("%s(%v) = %d, ожидается %d", function.name, tt.data, got, tt.want)
					}
					if !slices.Equal(data, tt.data) {
						t.Errorf("%s изменила исходные данные: %v, было %v", function.name, data, tt.data)
					}
				})
			}
		})
	}
}

func TestFindMax(t *testing.T) {
	tests := []struct {
		name        string
		data        []int
		start, stop int
		want        int
	}{
		{name: "весь диапазон", data: []int{3, 8, 2}, start: 0, stop: 3, want: 8},
		{name: "максимумы вне диапазона", data: []int{99, -5, -2, 100}, start: 1, stop: 3, want: -2},
		{name: "последний элемент", data: []int{9, 2}, start: 1, stop: 2, want: 2},
		{name: "правая граница исключена", data: []int{2, 9}, start: 0, stop: 1, want: 2},
		{name: "nil", data: nil, start: 0, stop: 0, want: 0},
		{name: "пустой слайс", data: []int{}, start: 0, stop: 1, want: 0},
		{name: "пустой диапазон", data: []int{1, 2}, start: 1, stop: 1, want: 0},
		{name: "обратные границы", data: []int{1, 2}, start: 2, stop: 1, want: 0},
		{name: "отрицательное начало", data: []int{1, 2}, start: -1, stop: 2, want: 0},
		{name: "отрицательный конец", data: []int{1, 2}, start: 0, stop: -1, want: 0},
		{name: "обе границы отрицательные", data: []int{1, 2}, start: -2, stop: -1, want: 0},
		{name: "конец за слайсом", data: []int{1, 2}, start: 0, stop: 3, want: 0},
		{name: "начало на длине", data: []int{1, 2}, start: 2, stop: 3, want: 0},
		{name: "начало за слайсом", data: []int{1, 2}, start: 3, stop: 4, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMax(tt.data, tt.start, tt.stop); got != tt.want {
				t.Errorf("findMax(%v, %d, %d) = %d, ожидается %d", tt.data, tt.start, tt.stop, got, tt.want)
			}
		})
	}
}

func TestApplyFuncToChunk(t *testing.T) {
	tests := []struct {
		name                  string
		size                  int
		nilWaitGroup, nilFunc bool
		want                  [][2]int
	}{
		{name: "nil WaitGroup", size: 8, nilWaitGroup: true},
		{name: "nil функция", size: 8, nilFunc: true},
		{name: "нулевой размер", size: 0},
		{name: "отрицательный размер", size: -1},
		{name: "один элемент", size: 1, want: [][2]int{{0, 1}}},
		{name: "меньше восьми элементов", size: 3, want: [][2]int{{0, 1}, {1, 2}, {2, 3}}},
		{name: "восемь элементов", size: 8, want: [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}}},
		{name: "делится на восемь", size: 16, want: [][2]int{{0, 2}, {2, 4}, {4, 6}, {6, 8}, {8, 10}, {10, 12}, {12, 14}, {14, 16}}},
		{name: "один лишний элемент", size: 17, want: [][2]int{{0, 3}, {3, 5}, {5, 7}, {7, 9}, {9, 11}, {11, 13}, {13, 15}, {15, 17}}},
		{name: "несколько лишних элементов", size: 19, want: [][2]int{{0, 3}, {3, 6}, {6, 9}, {9, 11}, {11, 13}, {13, 15}, {15, 17}, {17, 19}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			group := &wg
			if tt.nilWaitGroup {
				group = nil
			}
			var mu sync.Mutex
			got := make(map[int][2]int)
			calls := 0
			fn := func(wg *sync.WaitGroup, index, start, stop int) {
				defer wg.Done()
				mu.Lock()
				defer mu.Unlock()
				calls++
				got[index] = [2]int{start, stop}
			}
			if tt.nilFunc {
				fn = nil
			}
			applyFuncToChunk(group, tt.size, fn)
			mu.Lock()
			defer mu.Unlock()
			if calls != len(tt.want) || len(got) != len(tt.want) {
				t.Fatalf("число вызовов = %d, уникальных частей = %d, ожидается %d", calls, len(got), len(tt.want))
			}
			for i, want := range tt.want {
				if bounds, ok := got[i]; !ok || bounds != want {
					t.Errorf("часть %d: границы = %v (есть: %t), ожидается %v", i, bounds, ok, want)
				}
			}
		})
	}
}
