package tokenizer

import (
	"strings"

	"github.com/james-bowman/sparse"
	"github.com/mjboos/godedup/preprocessing"
)

// TODO: allow multiple analyzer types
// TODO: profile because it's much slower with a set

// what's a good implementation here?
// do it similarly to sklearn and extract the tokens?
// use the sparse package

type NGramVectorizer struct {
	Min_n, Max_n int
	Vocab        []string
}

func unique(words []string) []string {
	vocab := make(map[string]bool)
	for _, word := range words {
		if !vocab[word] {
			vocab[word] = true
		}
	}
	return preprocessing.GetKeys(vocab)
}

func NewNGramVectorizer(min_n, max_n int, corpus []string) *NGramVectorizer {
	ngrams := GetCorpusNGrams(corpus, min_n, max_n)
	nGramVec := NGramVectorizer{min_n, max_n, unique(ngrams)}
	return &nGramVec
}

func (nGram *NGramVectorizer) Fit(corpus []string) *NGramVectorizer {
	ngrams := GetCorpusNGrams(corpus, nGram.Min_n, nGram.Max_n)
	nGram.Vocab = unique(ngrams)
	return nGram
}

func (nGram *NGramVectorizer) Transform(corpus []string) *sparse.DOK {
	dok := sparse.NewDOK(len(corpus), len(nGram.Vocab))
	for i, part := range corpus {
		for j, ngram := range nGram.Vocab {
			dok.Set(i, j, float64(strings.Count(part, ngram)))
		}
	}
	return dok
}

func preProcessWord(word string) string {
	return strings.Trim(strings.ToLower(word), " .,")
}

func getGrams(word string, gramLen int) []string {
	embed := strings.Repeat(" ", gramLen-1)
	embeddedWord := embed + word + embed
	ngrams := make([]string, len(embeddedWord)-gramLen+1)
	for i := 0; i < len(embeddedWord)-gramLen+1; i++ {
		ngrams[i] = embeddedWord[i : i+gramLen]
	}
	return ngrams
}

func GetCorpusNGrams(corpus []string, min_n, max_n int) []string {
	var ngrams []string
	for _, word := range corpus {
		ngrams = append(ngrams, MakeNGrams(word, min_n, max_n)...)
	}
	return ngrams
}

func MakeNGrams(word string, min_n, max_n int) []string {
	preppedWord := preProcessWord(word)
	var grams []string
	for i := min_n; i <= max_n; i++ {
		grams = append(grams, getGrams(preppedWord, i)...)
	}
	return grams
}

func NGramDistance(s1 string, s2 string, lower int, upper int) float64 {
	counts := 0
	for _, w := range s1 {
		if strings.Contains(s2, string(w)) {
			counts++
		}
	}
	return float64(counts) / float64(len(s1))
}

func MakeNGramDistanceFunc(lower int, upper int) func(string, string) float64 {
	return func(s1 string, s2 string) float64 {
		return NGramDistance(s1, s2, lower, upper)
	}
}
