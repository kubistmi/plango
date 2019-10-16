// Package utils ...
package utils

import "fmt"

// FindMin ...
func FindMin(vec []int) int {
	var min int
	for ix, val := range vec {
		if ix == 0 || val < min {
			min = val
		}
	}
	return min
}

// FindMax ...
func FindMax(vec []int) int {
	var max int
	for ix, val := range vec {
		if ix == 0 || val > max {
			max = val
		}
	}
	return max
}

// MakeRange ...
func MakeRange(min, max int) ([]int, error) {
	if min >= max {
		return []int{}, fmt.Errorf("The parameter `max` must be strictly larger than parameter `min`")
	}

	vec := make([]int, max-min+1)
	for ix := range vec {
		vec[ix] = min + ix
	}
	return vec, nil
}

// FindUnique ...
func FindUnique(vec []int) []int {
	unique := make([]int, 0, len(vec))
	mapper := make(map[int]bool)

	for _, val := range vec {
		if _, ok := mapper[val]; !ok {
			mapper[val] = true
			unique = append(unique, val)
		}
	}
	return unique
}
