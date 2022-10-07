package preprocessing

import "strings"

func WordCount(s string) map[string]int {
	uniqueCount := make(map[string]int)
	for _, word := range strings.Fields(s) {
		uniqueCount[word] += 1
	}
	return uniqueCount
}

func CreateVectorFormDist(words []string, distance func(string, string) float64) []float64 {
	wlen := len(words)
	distances := make([]float64, int(wlen*(wlen-1)/2))
	idx := 0
	for i, word := range words {
		for _, cmpWord := range words[(i + 1):] {
			distances[idx] = distance(word, cmpWord)
			idx++
		}
	}
	return distances
}

func GetKeys[K comparable, V any](mapping map[K]V) []K {
	keys := make([]K, len(mapping))
	i := 0
	for k := range mapping {
		keys[i] = k
		i++
	}
	return keys
}
