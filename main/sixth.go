package main

import (
	"fmt"
	"math/rand"
	"time"
)
var r = rand.New(rand.NewSource(time.Now().UnixNano()))
func CreateRandomGenerator(lowerBound int, upperBound int) int {
	
	return r.Intn(upperBound-lowerBound+1) + lowerBound
}

func ForLoopExample(x int) {
	for i := 0; i < x; i += 2 {
		if i > 10 {
			fmt.Println("I am greater than 10")
			break
		}
	}
}

func CreateAndFilledArray(array []int) []int {
	
	var arr []int = make([]int, len(array))
	for i := 0; i < len(array); i++ {
		arr[i] = CreateRandomGenerator(1, 100)
	}

	for i := 0; i < len(arr); i++ {
		if array[i] < arr[i] {
			array[i] = arr[i]
		}
	}
	return arr
}

