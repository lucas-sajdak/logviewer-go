package main

// Min returns minimum of two integers
func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Max returns maximum of two integers
func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
