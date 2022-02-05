package util

// GenerateContinuousSlice
// input: low, high uint64
// output: [low, high] including the left and right boundaries
func GenerateContinuousSlice(low, high uint64) []uint64 {
    res := make([]uint64, high-low+1)
    for i := low; i <= high; i++ {
    	res[i-low] = i
	}
	return res
}
