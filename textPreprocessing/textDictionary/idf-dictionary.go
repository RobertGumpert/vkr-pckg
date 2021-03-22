package textDictionary

import (
	"github.com/RobertGumpert/vkr-pckg/textPreprocessing"
	concurrentMap "github.com/streamrail/concurrent-map"
	"sync"
)

//
//
//
//-------------------------IDF DICTIONARY-------------------------------------------------------------------------------
//
//
//

func IDFDictionary(corpus []string, idf int64, mode textPreprocessing.ThreadMode) (concurrentMap.ConcurrentMap, [][]string, int) {
	if mode == textPreprocessing.ParallelMode {
		return parallelCreateIDFDictionary(corpus, idf)
	}
	return linearCreateIDFDictionary(corpus, idf)
}

func linearCreateIDFDictionary(corpus []string, idf int64) (concurrentMap.ConcurrentMap, [][]string, int) {
	var (
		dictionary     = concurrentMap.New()
		vectorsOfWords = make([][]string, len(corpus))
	)
	for text := 0; text < len(corpus); text++ {
		if corpus[text] == "" || corpus[text] == "\n" || corpus[text] == "\r" || corpus[text] == "\r\n" {
			continue
		}
		addWordsToIDFDictionary(dictionary, corpus[text], text, vectorsOfWords)
	}
	dictionaryTransformIDF(dictionary, idf)
	return dictionary, vectorsOfWords, dictionary.Count()
}

func parallelCreateIDFDictionary(corpus []string, idf int64) (concurrentMap.ConcurrentMap, [][]string, int) {
	var (
		wg             = new(sync.WaitGroup)
		dictionary     = concurrentMap.New()
		vectorsOfWords = make([][]string, len(corpus))
	)
	for text := 0; text < len(corpus); text++ {
		if corpus[text] == "" || corpus[text] == "\n" || corpus[text] == "\r" || corpus[text] == "\r\n" {
			continue
		}
		wg.Add(1)
		go func(dictionary concurrentMap.ConcurrentMap, text string, positionTextInCorpus int, vectorsOfWords [][]string, wg *sync.WaitGroup) {
			defer wg.Done()
			addWordsToIDFDictionary(dictionary, text, positionTextInCorpus, vectorsOfWords)
			return
		}(dictionary, corpus[text], text, vectorsOfWords, wg)
	}
	wg.Wait()
	dictionaryTransformIDF(dictionary, idf)
	return dictionary, vectorsOfWords, dictionary.Count()
}

func addWordsToIDFDictionary(dictionary concurrentMap.ConcurrentMap, text string, positionInVectorsOfWords int, vectorsOfWords [][]string) {
	var (
		buffer = concurrentMap.New()
	)
	words := TextTransformToFeaturesSlice(text)
	if len(words) == 0 {
		return
	}
	for word := 0; word < len(words); word++ {
		vectorsOfWords[positionInVectorsOfWords] = append(
			vectorsOfWords[positionInVectorsOfWords],
			words[word],
		)
		existInBuffer := buffer.Has(words[word])
		itemInDictionary, existInDictionary := dictionary.Get(words[word])
		if !existInBuffer {
			buffer.Set(words[word], struct{}{})
		}
		if existInDictionary && !existInBuffer {
			freq := itemInDictionary.(int64) + 1
			dictionary.Upsert(words[word], freq, func(exist bool, valueInMap interface{}, newValue interface{}) interface{} {
				return newValue
			})
		}
		if !existInDictionary {
			dictionary.Set(words[word], int64(1))
		}
	}
}

func dictionaryTransformIDF(dictionary concurrentMap.ConcurrentMap, idf int64) {
	words := dictionary.Keys()
	index := 0
	for position := 0; position < len(words); position++ {
		item, _ := dictionary.Get(words[position])
		if item.(int64) < idf {
			dictionary.Remove(words[position])
			continue
		}
		dictionary.Upsert(words[position], int64(index), func(exist bool, valueInMap interface{}, newValue interface{}) interface{} {
			return newValue
		})
		index++
	}
}
