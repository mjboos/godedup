package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mjboos/godedup/preprocessing"
	"github.com/mjboos/godedup/tokenizer"
	//	"github.com/knightjdr/hclust"
)

func IsUpper(to_test string) bool {
	if strings.ToUpper(to_test) == to_test {
		return true
	} else {
		return false
	}
}

func GetIndices(line string) []int {
	var splits []int
	for idx := 0; idx < len(line); idx += 1 {
		if IsUpper(string(line[idx])) {
			splits = append(splits, idx)
		}
	}
	return splits
}

func SplitLine(line string) []string {
	splits := GetIndices(line)
	var lines []string
	for idx, from := range splits {
		var to int
		if next_idx := idx + 1; next_idx >= len(splits) {
			to = len(line) - 1
		} else {
			to = splits[next_idx]
		}
		lines = append(lines, line[from:to])
	}
	return lines
}

func UnRollLines(lines []string) []string {
	var corrLines []string
	for _, line := range lines {
		new_lines := SplitLine(line)
		corrLines = append(corrLines, new_lines...)
	}
	return corrLines
}

func main() {
	readFile, err := os.Open("example_data.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		txt := fileScanner.Text()
		fileLines = append(fileLines, txt)
	}
	readFile.Close()
	corrLines := UnRollLines(fileLines)
	uniqueWords := preprocessing.WordCount(strings.Join(corrLines, string(' ')))
	fmt.Println(uniqueWords)
	keys := preprocessing.GetKeys(uniqueWords)
	vec := tokenizer.NewNGramVectorizer(2, 3, keys)
	fmt.Println(vec.Vocab)
	mat := vec.Transform([]string{"zioziozio"})
	x, y := mat.Dims()
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			val := mat.At(i, j)
			if val > 0 {
				fmt.Println("######")
				fmt.Println(vec.Vocab[j])
				fmt.Println(val)
			}
		}
	}
	//	ngramDist := tokenizer.MakeNGramDistanceFunc(2, 3)
	//	vecDist := preprocessing.CreateVectorFormDist(keys, ngramDist)
	//	fmt.Println(vecDist)

}
