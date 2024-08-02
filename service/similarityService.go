package service

import (
	"github.com/go-nlp/bm25"
	"github.com/go-nlp/tfidf"
	"similarty-engine/model"
	"sort"
	"strings"
)

type SimilarityService struct {
}

func NewSimilarityService() SimilarityService {
	return SimilarityService{}
}

type doc []int

func (ss SimilarityService) FilterLines(lines []model.Line, frequency float64) model.FilterLines {
	library, corpus, documents := ss.prepareLibrary(lines)
	result := make([]model.Line, 0)

	ignoredLineIndex := make(map[int]bool)
	for i := 0; i < len(lines); i++ {
		str1 := lines[i].Title
		document := makeDocument(str1, corpus)

		if !ignoredLineIndex[i] {

			documentScores := bm25.BM25(library, doc(document), documents, 1.5, 0.5)
			sort.Sort(sort.Reverse(documentScores))

			topResult := selectTopPercentage(documentScores, frequency)
			for _, d := range documentScores {
				for _, ts := range topResult {
					if ts == d.Score {
						ignoredLineIndex[d.ID] = true
						break
					}
				}
			}

			result = append(result, lines[i])
		}

		//println("--------------------------")

	}

	return model.FilterLines{Lines: result}
}

func selectTopPercentage(dScores bm25.DocScores, percentage float64) []float64 {
	scores := make([]float64, 0)
	for _, s := range dScores {
		scores = append(scores, s.Score)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))

	var totalScore float64
	for _, score := range scores {
		totalScore += score
	}

	target := totalScore * percentage

	var selectedScores []float64
	var accumulatedScore float64
	for _, score := range scores {
		if accumulatedScore >= target {
			break
		}
		selectedScores = append(selectedScores, score)
		accumulatedScore += score
	}

	return selectedScores
}

func (ss SimilarityService) prepareLibrary(lines []model.Line) (*tfidf.TFIDF, map[string]int, []tfidf.Document) {
	corpus, _ := makeCorpus(lines)
	docs := makeDocuments(lines, corpus)
	tf := tfidf.New()

	for _, doc := range docs {
		tf.Add(doc)
	}
	tf.CalculateIDF()

	return tf, corpus, docs
}

func (d doc) IDs() []int { return d }

func makeCorpus(a []model.Line) (map[string]int, []string) {
	retVal := make(map[string]int)
	invRetVal := make([]string, 0)
	var id int
	for _, s := range a {
		for _, f := range strings.Fields(s.Title) {
			f = strings.ToLower(f)
			if _, ok := retVal[f]; !ok {
				retVal[f] = id
				invRetVal = append(invRetVal, f)
				id++
			}
		}
	}
	return retVal, invRetVal
}

func makeDocuments(a []model.Line, c map[string]int) []tfidf.Document {
	retVal := make([]tfidf.Document, 0, len(a))
	for _, s := range a {
		retVal = append(retVal, doc(makeDocument(s.Title, c)))
	}
	return retVal
}

func makeDocument(a string, c map[string]int) []int {
	var ts []int
	for _, f := range strings.Fields(a) {
		f = strings.ToLower(f)
		id := c[f]
		ts = append(ts, id)
	}

	return ts
}
