package util

// assuming target > 0
func FloorPowerOf2(target uint64) uint64 {
	var maximum uint64
	maximum = 1 << 63
	if target >= maximum {
		return maximum
	}
	target = target >> 1
	target |= target >> 1
	target |= target >> 2
	target |= target >> 4
	target |= target >> 8
	target |= target >> 16
	target |= target >> 32

	return target + 1
}
