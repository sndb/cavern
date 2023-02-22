package main

import "math/rand"

// randPop removes and returns a random element from the list and
// the updated list.
func randPop[T any](slice []T) (T, []T) {
	i := rand.Intn(len(slice))
	elem := slice[i]
	slice[i] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return elem, slice
}

// in checks if elem is in slice.
func in[T comparable](slice []T, elem T) bool {
	for _, x := range slice {
		if x == elem {
			return true
		}
	}
	return false
}

// abs returns the absolute value of n.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
