package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

func Top10(input string) []string {
	wordCountsMap := make(map[string]int)

	for _, word := range strings.Fields(input) {
		wordCountsMap[word]++
	}

	wordCounts := make([]wordCount, 0, len(wordCountsMap))

	for word, count := range wordCountsMap {
		wordCounts = append(wordCounts, wordCount{word, count})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].count == wordCounts[j].count {
			return wordCounts[i].word < wordCounts[j].word
		}

		return wordCounts[i].count > wordCounts[j].count
	})

	var topWords []wordCount

	if len(wordCounts) >= 10 {
		topWords = wordCounts[:10]
	} else {
		topWords = wordCounts
	}

	result := make([]string, 0, len(topWords))

	for _, wc := range topWords {
		result = append(result, wc.word)
	}

	return result
}
