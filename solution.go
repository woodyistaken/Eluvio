// without parallelizing, we get a 13.5-15 ms
// 56 13312
package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type commonBitstring struct {
	currentLength int
	currentIndex  int
	bestLength    int
	bestIndex     int
	lock          *sync.Mutex
}

func main() {
	start := time.Now()
	bestIndex, bestLength := findLongestCommonSubstring()
	elapsed := time.Since(start) // compute duration
	fmt.Println("Execution time:", elapsed)
	println(bestIndex)
	println(bestLength)
}
func findLongestCommonSubstring() (int, int) {

	files := [10]string{"sample.1", "sample.2", "sample.3", "sample.4", "sample.5", "sample.6", "sample.7", "sample.8", "sample.9", "sample.10"}
	//files := [2]string{"sample.1", "sample.2"}

	fileData := make([][]byte, len(files))
	for i := 0; i < len(files); i++ {
		loadFiles(files[i], &fileData, i)
	}
	var bitstringMapLock sync.Mutex
	bitstringMap := make(map[int]*commonBitstring) //12:bitstringInfo mappings
	finishedFiles := make([]bool, len(files))
	finishLoop := false
	maxFileLength := maxLength(&fileData)

	var wg sync.WaitGroup

	for index := 0; index < maxFileLength; index++ { //Loop through the files

		finishLoop = true //Checking if every files has been checked
		for i := 0; i < len(finishedFiles); i++ {
			if finishedFiles[i] == false {
				finishLoop = false
				break
			}
		}
		if finishLoop {
			break
		}
		//wg.Add(1)
		checkIndexByte(index, &bitstringMap, fileData, &finishedFiles, &bitstringMapLock, &wg)

	}
	//wg.Wait()
	bestCommonBitstring := commonBitstring{currentLength: 0, currentIndex: 0, bestLength: 0, bestIndex: 0}
	bestIndex := -1
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if bestCommonBitstring.bestLength < bitstringMap[i*10+j].bestLength {
				bestCommonBitstring = *bitstringMap[i*10+j]
				bestIndex = i*10 + j
			}
		}
	}
	return bestIndex, bestCommonBitstring.bestLength
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

func maxLength(arrayOfArray *[][]byte) int {
	maxLen := 0
	for _, arr := range *arrayOfArray {
		if len(arr) > maxLen {
			maxLen = len(arr)
		}
	}
	return maxLen
}

func checkIndexByte(index int, bitstringMap *map[int]*commonBitstring, fileData [][]byte, finishedFiles *[]bool, bitstringMapLock *sync.Mutex, wg *sync.WaitGroup) {
	//defer wg.Done()
	for i := 0; i < len(fileData); i++ { // Comparing each file bit to the other files
		if index >= len(fileData[i]) { //Checking if fileData[i] has been finished
			(*finishedFiles)[i] = true
			continue
		}

		for j := i + 1; j < len(fileData); j++ {
			if index >= len(fileData[j]) { //Checking if fileData[j] has been finished
				(*finishedFiles)[j] = true
				continue
			}
			compareBytes(i, j, index, bitstringMap, fileData, bitstringMapLock)
		}

	}
}

func compareBytes(indexA int, indexB int, index int, bitstringMap *map[int]*commonBitstring, fileData [][]byte, mapLock *sync.Mutex) {

	mapIndex := indexA*10 + indexB
	dataAByte := fileData[indexA][index]
	dataBByte := fileData[indexB][index]

	//mapLock.Lock()
	bitstringInfo, exists := (*bitstringMap)[mapIndex]
	if !exists {
		(*bitstringMap)[mapIndex] = &commonBitstring{currentLength: 0, currentIndex: 0, bestLength: 0, bestIndex: 0, lock: &sync.Mutex{}}
		bitstringInfo = (*bitstringMap)[mapIndex]
	}
	//mapLock.Unlock()

	//bitstringInfo.lock.Lock()
	if dataAByte == dataBByte {
		bitstringInfo.currentLength += 1
	} else {
		if bitstringInfo.currentLength > bitstringInfo.bestLength {
			bitstringInfo.bestLength = bitstringInfo.currentLength
			bitstringInfo.bestIndex = bitstringInfo.currentIndex
		}
		bitstringInfo.currentLength = 0
		bitstringInfo.currentIndex = index + 1
	}
	//bitstringInfo.lock.Unlock()
}
