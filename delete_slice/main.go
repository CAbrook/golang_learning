package main

import (
	"fmt"
)

func main() {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	PrintSlice(&slice)
	removeSliceEle(2, &slice)
	PrintSlice(&slice)
	removeSliceEle(2, &slice)
	removeSliceEle(2, &slice)
	removeSliceEle(2, &slice)
	removeSliceEle(2, &slice)
	removeSliceEle(2, &slice)
	PrintSlice(&slice)
}

func removeSliceEle[T any](index int, slice *[]T) {
	if index > len(*slice)-1 {
		panic("index out of range")
	}
	copy((*slice)[index:], (*slice)[index+1:])
	var zero T
	(*slice)[len(*slice)-1] = zero
	*slice = (*slice)[:len(*slice)-1]

	if len(*slice) <= (cap(*slice)+3)/4 {
		newSlice := make([]T, len(*slice), (cap(*slice)+1)/2)
		copy(newSlice, (*slice))
		*slice = newSlice
	}
}

func PrintSlice[T any](slice *[]T) {
	fmt.Printf("slice: %v, length: %d, capacity: %d \n", *slice, len(*slice), cap(*slice))
}
