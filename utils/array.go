package utils

// GetMostFrequentValue returns the most frequent value in the array.
// Example: [1, 2, 2, 3, 3, 3] -> 3, 3
// Example: [1, 2, 2, 3, 3, 3, 4, 4, 4, 4] -> 4, 4
// Example: [1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5] -> 5, 5
func GetMostFrequentValue(values []int) (int, int) {
	// get the most frequent value
	var (
		mostFrequentValue int
		maxCount          int
	)

	for _, value := range values {
		count := 0
		for _, v := range values {
			if v == value {
				count++
			}
		}
		if count > maxCount {
			maxCount = count
			mostFrequentValue = value
		}
	}
	return mostFrequentValue, maxCount
}
