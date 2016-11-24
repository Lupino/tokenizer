package main

import (
	"sync"

	"github.com/huichen/sego"
)

var (
	dictSegmenterMap      = map[string]*sego.Segmenter{}
	dictSegmenterMapMutex sync.Mutex
)

func getSegoSegmenter(dictFiles string) (*sego.Segmenter, error) {
	dictSegmenterMapMutex.Lock()
	defer dictSegmenterMapMutex.Unlock()

	if segmenter, ok := dictSegmenterMap[dictFiles]; ok {
		return segmenter, nil
	}

	segmenter := new(sego.Segmenter)
	segmenter.LoadDictionary(dictFiles)

	dictSegmenterMap[dictFiles] = segmenter

	return segmenter, nil
}
