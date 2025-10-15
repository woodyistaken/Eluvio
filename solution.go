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
)

// commonBitstring struct keeps track of the longest substring in a pair of files.
type commonBitstring struct {
	fileA      string //file name for A
	fileB      string //file name for B
	indexA     int    //index of longest common substring of fileA and fileB for fileA
	indexB     int    //index of longest common substring of fileA and fileB for fileB
	bestLength int    //length of the longest common substring between file A and file B
}

func main() {
	fileNames := os.Args[1:]
	bestLength, bestFileNames, bestOffsets := findLongestCommonSubstring(fileNames)

	println("length:", bestLength)
	for i := range bestOffsets {
		print(bestFileNames[i])
		print(" : ")
		println("start at", bestOffsets[i])
	}
}

// The function finds the longest common substring from the list of files names in the argument.
// Returns the length of the common substring, an array of filenames with the common substring, and the index of the common substring in those files
func findLongestCommonSubstring(files []string) (int, []string, []int) {

	//fileData is an array of file data from each file
	fileData := make([][]byte, len(files))
	for i := 0; i < len(files); i++ {
		loadFiles(files[i], &fileData, i)
	}

	bitstringMap := allocateWorkers(&fileData, &files)

	var bestCommonBitstrings []commonBitstring
	maxlen := 0
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			mapIndex := files[i] + "--" + files[j]
			if maxlen < bitstringMap[mapIndex].bestLength {
				bestCommonBitstrings = make([]commonBitstring, 0)
				bestCommonBitstrings = append(bestCommonBitstrings, *bitstringMap[mapIndex])
				maxlen = bitstringMap[mapIndex].bestLength
			}
			if maxlen == bitstringMap[mapIndex].bestLength {
				bestCommonBitstrings = append(bestCommonBitstrings, *bitstringMap[mapIndex])
			}
		}
	}

	var bestFileNames []string
	var bestIndices []int
	for i := range bestCommonBitstrings {
		fileA := bestCommonBitstrings[i].fileA
		fileB := bestCommonBitstrings[i].fileB
		indexA := bestCommonBitstrings[i].indexA
		indexB := bestCommonBitstrings[i].indexB
		if !strContains(bestFileNames, fileA) {
			bestFileNames = append(bestFileNames, fileA)
			bestIndices = append(bestIndices, indexA)

		}
		if !strContains(bestFileNames, fileB) {
			bestFileNames = append(bestFileNames, fileB)
			bestIndices = append(bestIndices, indexB)
		}
	}
	return maxlen, bestFileNames, bestIndices
}

// strContains takes in an array of strings and check if val is in the array
// Returns true if val is in the array, otherwise false
func strContains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// loadFiles take a filename and load the file data into fileData at index
func loadFiles(fileName string, fileData *[][]byte, index int) {
	data, err := os.ReadFile(fileName) // binary file
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	(*fileData)[index] = make([]byte, len(data))
	copy((*fileData)[index], data)
}

// checkIndexByte takes a pointer to array of fileData to compute the longest common substring between each pair of file
// Returns an array of pointer to commonBitstring struct that contains the length of longest common substring and the index the substring in each file.
func allocateWorkers(fileData *[][]byte, files *[]string) map[string]*commonBitstring {
	var bitstringMapLock sync.Mutex
	bitstringMap := make(map[string]*commonBitstring) //12:bitstringInfo mappings
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	jobs := make(chan [2]int, 100)
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pair := range jobs {
				i, j := pair[0], pair[1]
				longestSubstringBetweenTwoFiles(i, j, &bitstringMap, fileData, files, &bitstringMapLock)
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

// compareBytes takes data to compute the longest common substring between two files
/*
Parameter:
	indexA: index of fileA in the files array
	indexB: index of fileB in the files array
	bitstringMap: a map to keep track of longest common substring between two files
	fileData: fileData of all the files
	files: file names of all the files
	mapLock: synchronize read and write access to the bitstringMap
*/
func longestSubstringBetweenTwoFiles(indexA int, indexB int, bitstringMap *map[string]*commonBitstring, fileData *[][]byte, files *[]string, mapLock *sync.Mutex) {
	dataA := (*fileData)[indexA]
	dataB := (*fileData)[indexB]
	lengthA := len(dataA)
	lengthB := len(dataB)
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

	mapIndex := (*files)[indexA] + "--" + (*files)[indexB]
	mapLock.Lock()
	bitstringInfo, exists := (*bitstringMap)[mapIndex]
	if !exists {
		(*bitstringMap)[mapIndex] = &commonBitstring{fileA: (*files)[indexA], fileB: (*files)[indexB], indexA: 0, indexB: 0, bestLength: 0}
		bitstringInfo = (*bitstringMap)[mapIndex]
	}
	mapLock.Unlock()
	bitstringInfo.indexA = endAIndex - maxLen
	bitstringInfo.indexB = endBIndex - maxLen
	bitstringInfo.bestLength = maxLen
}
