// without parallelizing, we get a 35.7570769s single core
// with parallelizing, we get 4.28 seconds
// length 27648
// sample.2 : 3072
// sample.3 : 17408
package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type commonBitstring struct {
	mapIndex   int
	indexA     int
	indexB     int
	bestLength int
}

func main() {
	fileNames := os.Args[1:]
	var totalTime time.Duration
	start := time.Now()
	bestLength, bestFileNames, bestOffsets := findLongestCommonSubstring(fileNames)
	totalTime += time.Since(start)

	fmt.Println("Execution time:", totalTime)

	println("length:", bestLength)
	for i := range bestOffsets {
		print(bestFileNames[i])
		print(" : ")
		println("start at", bestOffsets[i])
	}
}
func findLongestCommonSubstring(files []string) (int, []string, []int) {

	//files := [10]string{"sample.1", "sample.2", "sample.3", "sample.4", "sample.5", "sample.6", "sample.7", "sample.8", "sample.9", "sample.10"}
	//files := [5]string{"pain.1", "pain.2", "pain.3", "pain.4", "pain.5"}
	//files := [3]string{"simple.1", "simple.2", "simple.3"}
	//files := [4]string{"sample.1", "sample.2", "sample.3", "sample.4"}

	fileData := make([][]byte, len(files))
	for i := 0; i < len(files); i++ {
		loadFiles(files[i], &fileData, i)
	}

	bitstringMap := checkIndexByte(&fileData)

	var bestCommonBitstrings []commonBitstring
	maxlen := 0
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if maxlen < bitstringMap[i*10+j].bestLength {
				bestCommonBitstrings = make([]commonBitstring, 0)
				bestCommonBitstrings = append(bestCommonBitstrings, *bitstringMap[i*10+j])
				maxlen = bitstringMap[i*10+j].bestLength
			}
			if maxlen == bitstringMap[i*10+j].bestLength {
				bestCommonBitstrings = append(bestCommonBitstrings, *bitstringMap[i*10+j])
			}
		}
	}

	var bestFileNames []string
	var bestIndices []int
	for i := range bestCommonBitstrings {
		fileA := int(bestCommonBitstrings[i].mapIndex / 10)
		fileB := int(bestCommonBitstrings[i].mapIndex % 10)
		indexA := bestCommonBitstrings[i].indexA
		indexB := bestCommonBitstrings[i].indexB
		if !strContains(bestFileNames, files[fileA]) {
			bestFileNames = append(bestFileNames, files[fileA])
			bestIndices = append(bestIndices, indexA)

		}
		if !strContains(bestFileNames, files[fileB]) {
			bestFileNames = append(bestFileNames, files[fileB])
			bestIndices = append(bestIndices, indexB)
		}
	}
	return maxlen, bestFileNames, bestIndices
}

func strContains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func loadFiles(fileName string, fileData *[][]byte, index int) {
	data, err := os.ReadFile(fileName) // binary file
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	(*fileData)[index] = make([]byte, len(data))
	copy((*fileData)[index], data)
}

func checkIndexByte(fileData *[][]byte) map[int]*commonBitstring {
	var bitstringMapLock sync.Mutex
	bitstringMap := make(map[int]*commonBitstring) //12:bitstringInfo mappings
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	jobs := make(chan [2]int, 100)
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pair := range jobs {
				i, j := pair[0], pair[1]
				compareBytes(i, j, (*fileData)[i], (*fileData)[j], len((*fileData)[i]), len((*fileData)[j]), &bitstringMap, fileData, &bitstringMapLock)
			}
		}()
	}
	for i := 0; i < len(*fileData); i++ { // Comparing each file bit to the other files
		for j := i + 1; j < len(*fileData); j++ {
			jobs <- [2]int{i, j}
		}
	}
	close(jobs)
	wg.Wait()
	return bitstringMap
}

func compareBytes(indexA int, indexB int, dataA []byte, dataB []byte, lengthA int, lengthB int, bitstringMap *map[int]*commonBitstring, fileData *[][]byte, mapLock *sync.Mutex) {

	prev := make([]int, lengthB+1)
	current := make([]int, lengthB+1)
	maxLen := 0
	endAIndex := 0
	endBIndex := 0
	for i := 1; i <= lengthA; i++ {
		for j := 1; j <= lengthB; j++ {
			if dataA[i-1] == dataB[j-1] {
				current[j] = prev[j-1] + 1
				if current[j] > maxLen {
					maxLen = current[j]
					endAIndex = i
					endBIndex = j
				}
			} else {
				current[j] = 0
			}
		}
		prev, current = current, prev
	}
	mapIndex := indexA*10 + indexB
	mapLock.Lock()
	bitstringInfo, exists := (*bitstringMap)[mapIndex]
	if !exists {
		(*bitstringMap)[mapIndex] = &commonBitstring{mapIndex: mapIndex, indexA: 0, indexB: 0, bestLength: 0}
		bitstringInfo = (*bitstringMap)[mapIndex]
	}
	mapLock.Unlock()
	bitstringInfo.indexA = endAIndex - maxLen
	bitstringInfo.indexB = endBIndex - maxLen
	bitstringInfo.bestLength = maxLen
}
