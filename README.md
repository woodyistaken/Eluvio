# Problem

Given a large number of binary files, write a program that finds the
longest strand of bytes that is identical between two or more files

Use the test set attached (files sample.*)

The program should display:
- the length of the strand
- the file names where the largest strand appears
- the offset where the strand appears in each file

# Solution

## How it works
1. **Loads all binary files** into memory  
2. **Compares each pair of files** using a dynamic programming approach to find the longest common substring between the pair
3. Uses **parallelization** for these comparisons to improve performance
4. Returns:
    - The **length** of the longest common substring
    - The **files** that contain the substring
    - The **offsets** where the substring starts

## Comparison Algorithm
For two files A and B
- Compute the longest substring endings at (i, j) for every pair of indices between the two files
- Stores two rows of indices to reduce memory usage from O(n^2)->O(n)
- Records the length of the longest common substring and starting index in each file

## Features

- Compares **multiple binary files** simultaneously  
- Uses **parallelization** via goroutines and worker pools for performance  
- Works directly with **binary data**  
- Thread-safe using sync.Mutex  
- Memory-efficient (uses rolling 2-row DP optimization)

## Performance


| Mode | Time (seconds) | CPU Cores Used |
|------|----------------|----------------|
| Single-core | 35.75s | 1 |
| Parallelized | **4.28s** | 16 |


## Example Usage
```bash
# Run program on the test set
go run solution.go "sample.1" "sample.2" "sample.3" "sample.4" "sample.5" "sample.6" "sample.7" "sample.8" "sample.9" "sample.10"
