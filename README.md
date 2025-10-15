# Problem

Given a large number of binary files, write a program that finds the
longest strand of bytes that is identical between two or more files

Use the test set attached (files sample.*)

The program should display:
- the length of the strand
- the file names where the largest strand appears
- the offset where the strand appears in each file

# Solution


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
| Parallelized | **4.28s** | All available cores |


## Example Usage
```bash
# Run program on the test set
go run solution.go "sample.1" "sample.2" "sample.3" "sample.4" "sample.5" "sample.6" "sample.7" "sample.8" "sample.9" "sample.10"
